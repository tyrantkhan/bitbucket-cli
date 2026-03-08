package pr

import (
	"context"
	"fmt"
	"os/exec"
	"sort"
	"strconv"
	"strings"

	"github.com/tyrantkhan/bb/internal/api"
	"github.com/tyrantkhan/bb/internal/cmdutil"
	"github.com/tyrantkhan/bb/internal/models"
	"github.com/tyrantkhan/bb/internal/output"
	"github.com/urfave/cli/v3"
)

func newCmdView() *cli.Command {
	return &cli.Command{
		Name:      "view",
		Usage:     "View pull request details",
		ArgsUsage: "<id>",
		Flags: []cli.Flag{
			cmdutil.WorkspaceFlag,
			cmdutil.RepoFlag,
			cmdutil.FormatFlag,
			cmdutil.WebFlag,
			&cli.BoolFlag{
				Name:  "comments",
				Usage: "Show PR comments",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			f := cmdutil.GetFactory(ctx)
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			workspace, repo, err := cmdutil.ResolveWorkspaceAndRepo(ctx, cmd)
			if err != nil {
				return err
			}

			idStr := cmd.Args().First()
			if idStr == "" {
				return fmt.Errorf("pull request ID is required")
			}
			prID, err := strconv.Atoi(idStr)
			if err != nil {
				return fmt.Errorf("invalid pull request ID: %s", idStr)
			}

			path := fmt.Sprintf("/2.0/repositories/%s/%s/pullrequests/%d", workspace, repo, prID)

			resp, err := client.Get(path)
			if err != nil {
				return err
			}

			var pr models.PullRequest
			if err := api.DecodeJSON(resp, &pr); err != nil {
				return fmt.Errorf("failed to decode pull request: %w", err)
			}

			// Open in browser if --web flag is set.
			if cmd.Bool("web") {
				url := ""
				if pr.Links.HTML != nil {
					url = pr.Links.HTML.Href
				}
				if url == "" {
					return fmt.Errorf("no web URL available for this pull request")
				}
				fmt.Fprintln(f.IOOut, output.Muted.Render("Opening in browser: "+url))
				return exec.Command("open", url).Run()
			}

			showComments := cmd.Bool("comments")
			format := cmdutil.GetFormat(ctx, cmd)

			if format == "json" {
				if showComments {
					comments, err := fetchComments(client, workspace, repo, prID)
					if err != nil {
						return err
					}
					return output.RenderJSON(struct {
						models.PullRequest
						Comments []models.Comment `json:"comments"`
					}{pr, comments})
				}
				return output.RenderJSON(pr)
			}

			// Rich detail view.
			fmt.Fprintf(f.IOOut, "%s #%d\n", output.Header.Render(pr.Title), pr.ID)
			fmt.Fprintf(f.IOOut, "%s  %s\n", output.Bold.Render("State:"), output.StatusColor(pr.State).Render(pr.State))
			if pr.Draft {
				fmt.Fprintf(f.IOOut, "%s  %s\n", output.Bold.Render("Draft:"), output.Yellow.Render("DRAFT"))
			}
			fmt.Fprintf(f.IOOut, "%s %s\n", output.Bold.Render("Author:"), pr.Author.DisplayName)
			fmt.Fprintf(f.IOOut, "%s %s → %s\n",
				output.Bold.Render("Branch:"),
				output.Cyan.Render(pr.Source.Branch.Name),
				output.Cyan.Render(pr.Destination.Branch.Name),
			)

			// Description
			if pr.Description != "" {
				fmt.Fprintln(f.IOOut)
				output.RenderMarkdown(pr.Description)
			}

			// Reviewers
			if len(pr.Participants) > 0 {
				fmt.Fprintln(f.IOOut)
				fmt.Fprintln(f.IOOut, output.Bold.Render("Reviewers:"))
				for _, p := range pr.Participants {
					if p.Role != "REVIEWER" {
						continue
					}
					status := "pending"
					statusStyle := output.Yellow
					if p.Approved {
						status = "approved"
						statusStyle = output.Green
					} else if p.State == "changes_requested" {
						status = "changes requested"
						statusStyle = output.Red
					}
					fmt.Fprintf(f.IOOut, "  %s  %s\n",
						p.User.DisplayName,
						statusStyle.Render(status),
					)
				}
			}

			// Counts
			fmt.Fprintln(f.IOOut)
			fmt.Fprintf(f.IOOut, "%s   %d\n", output.Bold.Render("Tasks:"), pr.TaskCount)
			fmt.Fprintf(f.IOOut, "%s %d\n", output.Bold.Render("Comments:"), pr.CommentCount)

			// Timestamps
			fmt.Fprintln(f.IOOut)
			fmt.Fprintf(f.IOOut, "%s  %s\n", output.Muted.Render("Created:"), models.FormatTime(pr.CreatedOn))
			fmt.Fprintf(f.IOOut, "%s  %s\n", output.Muted.Render("Updated:"), models.FormatTime(pr.UpdatedOn))

			if showComments {
				comments, err := fetchComments(client, workspace, repo, prID)
				if err != nil {
					return err
				}
				renderComments(f, comments)
			}

			return nil
		},
	}
}

func fetchComments(client *api.Client, workspace, repo string, prID int) ([]models.Comment, error) {
	path := fmt.Sprintf("/2.0/repositories/%s/%s/pullrequests/%d/comments", workspace, repo, prID)
	return api.Paginate[models.Comment](client, path, 0)
}

func renderComments(f *cmdutil.Factory, comments []models.Comment) {
	if len(comments) == 0 {
		fmt.Fprintln(f.IOOut)
		fmt.Fprintln(f.IOOut, output.Muted.Render("No comments."))
		return
	}

	type commentNode struct {
		comment  models.Comment
		children []*commentNode
	}

	nodeMap := make(map[int]*commentNode)
	var roots []*commentNode

	sort.Slice(comments, func(i, j int) bool {
		return comments[i].ID < comments[j].ID
	})

	for _, c := range comments {
		node := &commentNode{comment: c}
		nodeMap[c.ID] = node

		if c.Parent != nil {
			if parent, ok := nodeMap[c.Parent.ID]; ok {
				parent.children = append(parent.children, node)
			} else {
				roots = append(roots, node)
			}
		} else {
			roots = append(roots, node)
		}
	}

	fmt.Fprintln(f.IOOut)
	fmt.Fprintln(f.IOOut, output.Bold.Render("Comments:"))

	var renderNode func(node *commentNode, depth int)
	renderNode = func(node *commentNode, depth int) {
		indent := strings.Repeat("  ", depth)
		c := node.comment

		header := fmt.Sprintf("%s%s  %s",
			indent,
			output.Bold.Render(c.Author.DisplayName),
			output.Muted.Render(models.FormatTime(c.CreatedOn)),
		)
		fmt.Fprintln(f.IOOut, header)

		if c.Inline != nil {
			loc := indent + "  " + output.Cyan.Render(c.Inline.Path)
			if c.Inline.To != nil {
				loc += output.Cyan.Render(fmt.Sprintf(":%d", *c.Inline.To))
			}
			fmt.Fprintln(f.IOOut, loc)
		}

		if c.Deleted {
			fmt.Fprintln(f.IOOut, indent+"  "+output.Muted.Render("[deleted]"))
		} else {
			lines := strings.Split(c.Content.Raw, "\n")
			for _, line := range lines {
				fmt.Fprintln(f.IOOut, indent+"  "+line)
			}
		}

		fmt.Fprintln(f.IOOut)

		for _, child := range node.children {
			renderNode(child, depth+1)
		}
	}

	for _, root := range roots {
		renderNode(root, 0)
	}
}
