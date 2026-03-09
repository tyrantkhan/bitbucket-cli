package pr

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/tyrantkhan/bb/internal/api"
	"github.com/tyrantkhan/bb/internal/cmdutil"
	"github.com/tyrantkhan/bb/internal/models"
	"github.com/tyrantkhan/bb/internal/output"
	"github.com/urfave/cli/v3"
)

// ActivityUpdate represents a state-change update in a PR activity feed.
type ActivityUpdate struct {
	State  string      `json:"state"`
	Date   string      `json:"date"`
	Author models.User `json:"author"`
}

// ActivityApproval represents an approval event in a PR activity feed.
type ActivityApproval struct {
	Date string      `json:"date"`
	User models.User `json:"user"`
}

// ActivityItem represents a single activity entry which contains one of
// comment, update, or approval.
type ActivityItem struct {
	Comment  *models.Comment   `json:"comment,omitempty"`
	Update   *ActivityUpdate   `json:"update,omitempty"`
	Approval *ActivityApproval `json:"approval,omitempty"`
}

func newCmdActivity() *cli.Command {
	return &cli.Command{
		Name:      "activity",
		Usage:     "Show pull request activity",
		ArgsUsage: "<id>",
		Flags: []cli.Flag{
			cmdutil.WorkspaceFlag,
			cmdutil.RepoFlag,
			cmdutil.LimitFlag,
			cmdutil.FormatFlag,
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

			limit := int(cmd.Int("limit"))
			path := fmt.Sprintf("/2.0/repositories/%s/%s/pullrequests/%d/activity", workspace, repo, prID)

			rawItems, err := api.PaginateRaw(client, path, limit)
			if err != nil {
				return err
			}

			format := cmdutil.GetFormat(ctx, cmd)
			if format == "json" {
				var items []ActivityItem
				for _, raw := range rawItems {
					var item ActivityItem
					if err := json.Unmarshal(raw, &item); err != nil {
						fmt.Fprintf(f.IOErr, "warning: failed to unmarshal activity item: %v\n", err)
						continue
					}
					items = append(items, item)
				}
				return output.RenderJSON(items)
			}

			if len(rawItems) == 0 {
				fmt.Fprintln(f.IOOut, output.Muted.Render("No activity."))
				return nil
			}

			for _, raw := range rawItems {
				var item ActivityItem
				if err := json.Unmarshal(raw, &item); err != nil {
					continue
				}

				switch {
				case item.Approval != nil:
					a := item.Approval
					fmt.Fprintf(f.IOOut, "%s  %s  %s\n",
						output.Green.Render("APPROVED"),
						output.Bold.Render(a.User.DisplayName),
						output.Muted.Render(models.FormatTime(a.Date)),
					)

				case item.Update != nil:
					u := item.Update
					stateStyle := output.StatusColor(u.State)
					fmt.Fprintf(f.IOOut, "%s  %s  %s  %s\n",
						output.Yellow.Render("UPDATE"),
						output.Bold.Render(u.Author.DisplayName),
						stateStyle.Render(u.State),
						output.Muted.Render(models.FormatTime(u.Date)),
					)

				case item.Comment != nil:
					c := item.Comment
					fmt.Fprintf(f.IOOut, "%s  %s  %s\n",
						output.Blue.Render("COMMENT"),
						output.Bold.Render(c.Author.DisplayName),
						output.Muted.Render(models.FormatTime(c.CreatedOn)),
					)

					// Show a preview of the comment content.
					preview := c.Content.Raw
					if len(preview) > 120 {
						preview = preview[:120] + "..."
					}
					fmt.Fprintf(f.IOOut, "         %s\n", preview)
				}

				fmt.Fprintln(f.IOOut)
			}

			return nil
		},
	}
}
