package pr

import (
	"context"
	"fmt"
	"strconv"

	"github.com/tyrantkhan/bb/internal/api"
	"github.com/tyrantkhan/bb/internal/cmdutil"
	"github.com/tyrantkhan/bb/internal/models"
	"github.com/tyrantkhan/bb/internal/output"
	"github.com/urfave/cli/v3"
)

func newCmdComment() *cli.Command {
	return &cli.Command{
		Name:      "comment",
		Usage:     "Add a comment to a pull request",
		ArgsUsage: "<id>",
		Flags: []cli.Flag{
			cmdutil.WorkspaceFlag,
			cmdutil.RepoFlag,
			cmdutil.FormatFlag,
			&cli.StringFlag{
				Name:  "body",
				Usage: "Comment text",
			},
			&cli.StringFlag{
				Name:  "file",
				Usage: "File path for inline comment (used with --body)",
			},
			&cli.IntFlag{
				Name:  "line",
				Usage: "Line number for inline comment (used with --body and --file)",
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

			message := cmd.String("body")
			if message == "" {
				return fmt.Errorf("--body is required; use --body to provide comment text")
			}

			comment, err := addComment(client, workspace, repo, prID, cmd, message)
			if err != nil {
				return err
			}

			format := cmdutil.GetFormat(ctx, cmd)
			if format == "json" {
				return output.RenderJSON(comment)
			}

			fmt.Fprintln(f.IOOut, output.Success.Render(
				fmt.Sprintf("Comment #%d added to pull request #%d.", comment.ID, prID),
			))
			return nil
		},
	}
}

func addComment(client *api.Client, workspace, repo string, prID int, cmd *cli.Command, message string) (*models.Comment, error) {
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
		return nil, err
	}

	var comment models.Comment
	if err := api.DecodeJSON(resp, &comment); err != nil {
		return nil, fmt.Errorf("failed to decode comment: %w", err)
	}

	return &comment, nil
}
