package cmd

import (
	"context"
	"fmt"
	"os"
	"runtime/debug"
	"strings"

	authCmd "github.com/tyrantkhan/bb/cmd/auth"
	configCmd "github.com/tyrantkhan/bb/cmd/config"
	pipelineCmd "github.com/tyrantkhan/bb/cmd/pipeline"
	prCmd "github.com/tyrantkhan/bb/cmd/pr"
	repoCmd "github.com/tyrantkhan/bb/cmd/repo"
	searchCmd "github.com/tyrantkhan/bb/cmd/search"
	"github.com/tyrantkhan/bb/internal/cmdutil"
	"github.com/tyrantkhan/bb/internal/output"
	"github.com/tyrantkhan/bb/internal/update"
	"github.com/urfave/cli/v3"
)

// Build-time variables set via ldflags.
var (
	Version   = "dev"
	BuildDate = ""
)

func init() {
	if Version == "dev" {
		if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "(devel)" {
			Version = info.Main.Version
		} else if info != nil {
			// Built from source — try to extract VCS info for a dev version string.
			var revision, dirty string
			for _, s := range info.Settings {
				switch s.Key {
				case "vcs.revision":
					revision = s.Value
				case "vcs.modified":
					if s.Value == "true" {
						dirty = "-dirty"
					}
				}
			}
			if revision != "" {
				if len(revision) > 7 {
					revision = revision[:7]
				}
				Version = fmt.Sprintf("0.0.0-%s%s-dev", revision, dirty)
			}
		}
	}
}

func versionString() string {
	s := fmt.Sprintf("bb version %s", Version)
	if BuildDate != "" {
		s += fmt.Sprintf(" (%s)", BuildDate)
	}
	if Version != "dev" && !strings.HasSuffix(Version, "-dev") && strings.Contains(Version, ".") {
		s += fmt.Sprintf("\nhttps://github.com/tyrantkhan/bitbucket-cli/releases/tag/v%s", Version)
		s += "\nhttps://github.com/tyrantkhan/bitbucket-cli"
	}
	return s
}

// NewRootCommand creates the root bb command.
func NewRootCommand() *cli.Command {
	cli.VersionPrinter = func(cmd *cli.Command) {
		fmt.Println(versionString())
	}

	root := &cli.Command{
		Name:                  "bb",
		Usage:                 "Bitbucket Cloud CLI",
		Version:               Version,
		EnableShellCompletion: true,
		ConfigureShellCompletionCommand: func(cmd *cli.Command) {
			cmd.Hidden = false
			cmd.Usage = "Generate shell completion script (bash, zsh, fish, pwsh)"
		},
		CommandNotFound: cmdutil.CommandNotFound,
		Before: func(ctx context.Context, cmd *cli.Command) (context.Context, error) {
			f, err := cmdutil.NewFactory()
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", output.Error.Render(err.Error()))
				return ctx, err
			}
			return cmdutil.WithFactory(ctx, f), nil
		},
		After: func(ctx context.Context, cmd *cli.Command) error {
			if Version == "dev" || strings.HasSuffix(Version, "-dev") {
				return nil
			}
			if os.Getenv("BB_NO_UPDATE_NOTIFIER") == "1" {
				return nil
			}
			latest, err := update.CheckForUpdate(Version)
			if err != nil {
				return nil //nolint:nilerr // update check failures are non-fatal
			}
			if latest == "" {
				return nil
			}
			fmt.Fprintf(os.Stderr, "\n%s %s → %s\n%s %s\n",
				output.Yellow.Render("A new release of bb is available:"),
				output.Cyan.Render(Version),
				output.Cyan.Render(latest),
				output.Yellow.Render("To upgrade:"),
				output.Cyan.Render(update.UpgradeCommand()),
			)
			return nil
		},
		Commands: []*cli.Command{
			authCmd.NewCmdAuth(),
			configCmd.NewCmdConfig(),
			repoCmd.NewCmdRepo(),
			prCmd.NewCmdPR(),
			pipelineCmd.NewCmdPipeline(),
			searchCmd.NewCmdSearch(),
			{
				Name:    "version",
				Aliases: []string{"v"},
				Usage:   "Print the version",
				Hidden:  true,
				Action: func(ctx context.Context, cmd *cli.Command) error {
					fmt.Println(versionString())
					return nil
				},
			},
		},
	}

	root.Commands = append(root.Commands, newHelpTopicCommands(root)...)

	root.CustomRootCommandHelpTemplate = cli.RootCommandHelpTemplate +
		`HELP TOPICS:
` + helpTopicsSection() + "\n"

	return root
}
