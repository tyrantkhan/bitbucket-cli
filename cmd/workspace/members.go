package workspace

import (
	"context"
	"fmt"

	"github.com/tyrantkhan/bb/internal/api"
	"github.com/tyrantkhan/bb/internal/cmdutil"
	"github.com/tyrantkhan/bb/internal/models"
	"github.com/tyrantkhan/bb/internal/output"
	"github.com/urfave/cli/v3"
)

func newCmdMembers() *cli.Command {
	return &cli.Command{
		Name:  "members",
		Usage: "List members in a workspace",
		Flags: []cli.Flag{
			cmdutil.WorkspaceFlag,
			cmdutil.FormatFlag,
			cmdutil.LimitFlag,
		},
		Action: cmdutil.NoArgs(func(ctx context.Context, cmd *cli.Command) error {
			f := cmdutil.GetFactory(ctx)
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			workspace, err := cmdutil.ResolveWorkspace(ctx, cmd)
			if err != nil {
				return err
			}

			limit := int(cmd.Int("limit"))
			path := fmt.Sprintf("/2.0/workspaces/%s/members", workspace)

			members, err := api.Paginate[models.WorkspaceMembership](client, path, limit)
			if err != nil {
				return err
			}

			format := cmdutil.GetFormat(ctx, cmd)

			headers := []string{"Display Name", "Nickname", "Account ID"}
			rows := make([][]string, len(members))
			for i, m := range members {
				rows[i] = []string{
					m.User.DisplayName,
					m.User.Nickname,
					m.User.AccountID,
				}
			}

			return output.Format(format, members, headers, rows)
		}),
	}
}
