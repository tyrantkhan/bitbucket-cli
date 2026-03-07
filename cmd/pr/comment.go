package pr

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/tyrantkhan/bb/internal/api"
	"github.com/tyrantkhan/bb/internal/cmdutil"
	"github.com/tyrantkhan/bb/internal/models"
	"github.com/tyrantkhan/bb/internal/output"
	"github.com/urfave/cli/v3"
)

func newCmdComment() *cli.Command {
	return &cli.Command{
		Name:      "comment",
		Usage:     "List or add comments on a pull request",
		ArgsUsage: "<id>",
		Flags: []cli.Flag{
			cmdutil.WorkspaceFlag,
			cmdutil.RepoFlag,
			&cli.StringFlag{
				Name:  "add",
				Usage: "Add a new comment with the given message",
			},
			&cli.StringFlag{
				Name:  "file",
				Usage: "File path for inline comment (used with --add)",
			},
			&cli.IntFlag{
				Name:  "line",
				Usage: "Line number for inline comment (used with --add and --file)",
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

			message := cmd.String("add")

			if message != "" {
				return addComment(f, client, workspace, repo, prID, cmd, message)
			}

			return listComments(f, client, workspace, repo, prID)
		},
	}
}

func addComment(f *cmdutil.Factory, client *api.Client, workspace, repo string, prID int, cmd *cli.Command, message string) error {
	body := map[string]interface{}{
		"content": map[string]string{
			"raw": message,
		},
	}

	file := cmd.String("file")
	line := int(cmd.Int("line"))

	if file != "" {
		inline := map[string]interface{}{
			"path": file,
		}
		if line > 0 {
			inline["to"] = line
		}
		body["inline"] = inline
	}

	path := fmt.Sprintf("/2.0/repositories/%s/%s/pullrequests/%d/comments", workspace, repo, prID)

	resp, err := client.Post(path, body)
	if err != nil {
		return err
	}

	var comment models.Comment
	if err := api.DecodeJSON(resp, &comment); err != nil {
		return fmt.Errorf("failed to decode comment: %w", err)
	}

	fmt.Fprintln(f.IOOut, output.Success.Render(
		fmt.Sprintf("Comment #%d added to pull request #%d.", comment.ID, prID),
	))

	return nil
}

func listComments(f *cmdutil.Factory, client *api.Client, workspace, repo string, prID int) error {
	path := fmt.Sprintf("/2.0/repositories/%s/%s/pullrequests/%d/comments", workspace, repo, prID)

	comments, err := api.Paginate[models.Comment](client, path, 0)
	if err != nil {
		return err
	}

	if len(comments) == 0 {
		fmt.Fprintln(f.IOOut, output.Muted.Render("No comments."))
		return nil
	}

	// Build a thread tree: group replies by parent ID.
	type commentNode struct {
		comment  models.Comment
		children []*commentNode
	}

	nodeMap := make(map[int]*commentNode)
	var roots []*commentNode

	// Sort comments by ID for stable ordering.
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
				// Orphaned reply; treat as root.
				roots = append(roots, node)
			}
		} else {
			roots = append(roots, node)
		}
	}

	// Render the thread tree.
	var renderNode func(node *commentNode, depth int)
	renderNode = func(node *commentNode, depth int) {
		indent := strings.Repeat("  ", depth)
		c := node.comment

		// Header: author and timestamp.
		header := fmt.Sprintf("%s%s  %s",
			indent,
			output.Bold.Render(c.Author.DisplayName),
			output.Muted.Render(models.FormatTime(c.CreatedOn)),
		)
		fmt.Fprintln(f.IOOut, header)

		// Inline location if present.
		if c.Inline != nil {
			loc := indent + "  " + output.Cyan.Render(c.Inline.Path)
			if c.Inline.To != nil {
				loc += output.Cyan.Render(fmt.Sprintf(":%d", *c.Inline.To))
			}
			fmt.Fprintln(f.IOOut, loc)
		}

		// Content.
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

	return nil
}
