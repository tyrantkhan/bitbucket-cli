package pr

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"charm.land/huh/v2"
	"github.com/tyrantkhan/bb/internal/api"
	"github.com/tyrantkhan/bb/internal/cmdutil"
	"github.com/tyrantkhan/bb/internal/models"
	"github.com/tyrantkhan/bb/internal/output"
	"github.com/urfave/cli/v3"
)

func newCmdEdit() *cli.Command {
	return &cli.Command{
		Name:      "edit",
		Usage:     "Edit a pull request's title and description (use flags for reviewers)",
		ArgsUsage: "<id>",
		Flags: []cli.Flag{
			cmdutil.WorkspaceFlag,
			cmdutil.RepoFlag,
			cmdutil.FormatFlag,
			&cli.StringFlag{
				Name:  "title",
				Usage: "New title",
			},
			&cli.StringFlag{
				Name:  "description",
				Usage: "New description",
			},
			&cli.StringSliceFlag{
				Name:  "reviewer",
				Usage: "Replace all reviewers (UUID, repeatable)",
			},
			&cli.StringSliceFlag{
				Name:  "add-reviewer",
				Usage: "Add reviewer (UUID, repeatable)",
			},
			&cli.StringSliceFlag{
				Name:  "remove-reviewer",
				Usage: "Remove reviewer (UUID, repeatable)",
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

			title := pr.Title
			description := pr.Description

			hasFlags := cmd.IsSet("title") || cmd.IsSet("description") ||
				cmd.IsSet("reviewer") || cmd.IsSet("add-reviewer") || cmd.IsSet("remove-reviewer")

			if hasFlags {
				if cmd.IsSet("title") {
					title = cmd.String("title")
				}
				if cmd.IsSet("description") {
					description = cmd.String("description")
				}
			} else {
				// Interactive mode.
				err := huh.NewForm(
					huh.NewGroup(
						huh.NewInput().
							Title("Title").
							Value(&title),
						huh.NewText().
							Title("Description").
							Value(&description),
					),
				).Run()
				if err != nil {
					return err
				}
			}

			if title == "" {
				return fmt.Errorf("pull request title cannot be empty")
			}

			body := buildPRUpdateBody(pr)
			body["title"] = title
			body["description"] = description

			// Handle reviewer flags.
			if cmd.IsSet("reviewer") {
				uuids := cmd.StringSlice("reviewer")
				reviewers := make([]map[string]string, len(uuids))
				for i, uuid := range uuids {
					reviewers[i] = map[string]string{"uuid": uuid}
				}
				body["reviewers"] = reviewers
			} else if cmd.IsSet("add-reviewer") || cmd.IsSet("remove-reviewer") {
				existing := make(map[string]bool)
				for _, r := range pr.Reviewers {
					existing[r.UUID] = true
				}

				for _, uuid := range cmd.StringSlice("add-reviewer") {
					existing[uuid] = true
				}
				for _, uuid := range cmd.StringSlice("remove-reviewer") {
					delete(existing, uuid)
				}

				reviewers := make([]map[string]string, 0, len(existing))
				for uuid := range existing {
					reviewers = append(reviewers, map[string]string{"uuid": uuid})
				}
				body["reviewers"] = reviewers
			}

			resp, err = client.Put(path, body)
			if err != nil {
				return err
			}

			var updated models.PullRequest
			if err := api.DecodeJSON(resp, &updated); err != nil {
				return fmt.Errorf("failed to decode response: %w", err)
			}

			format := cmdutil.GetFormat(ctx, cmd)
			if format == "json" {
				return output.RenderJSON(updated)
			}

			fmt.Fprintln(f.IOOut, output.Success.Render(
				fmt.Sprintf("Pull request #%d updated.", prID),
			))

			var changes []string
			if updated.Title != pr.Title {
				changes = append(changes, fmt.Sprintf("Title: %s", updated.Title))
			}
			if updated.Description != pr.Description {
				changes = append(changes, "Description updated")
			}
			if cmd.IsSet("reviewer") || cmd.IsSet("add-reviewer") || cmd.IsSet("remove-reviewer") {
				names := make([]string, len(updated.Reviewers))
				for i, r := range updated.Reviewers {
					names[i] = r.DisplayName
				}
				if len(names) > 0 {
					changes = append(changes, fmt.Sprintf("Reviewers: %s", strings.Join(names, ", ")))
				} else {
					changes = append(changes, "Reviewers: none")
				}
			}

			for _, c := range changes {
				fmt.Fprintf(f.IOOut, "  %s\n", c)
			}

			return nil
		},
	}
}
