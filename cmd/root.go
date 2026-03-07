package cmd

import (
	"context"
	"fmt"
	"os"
	"runtime/debug"
	"strings"

	authCmd "github.com/tyrantkhan/bb/cmd/auth"
	pipelineCmd "github.com/tyrantkhan/bb/cmd/pipeline"
	prCmd "github.com/tyrantkhan/bb/cmd/pr"
	repoCmd "github.com/tyrantkhan/bb/cmd/repo"
	"github.com/tyrantkhan/bb/internal/cmdutil"
	"github.com/tyrantkhan/bb/internal/output"
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
		s += fmt.Sprintf("\nhttps://github.com/tyrantkhan/bb/releases/tag/v%s", Version)
	}
	return s
}

// NewRootCommand creates the root bb command.
func NewRootCommand() *cli.Command {
	cli.VersionPrinter = func(cmd *cli.Command) {
		fmt.Println(versionString())
	}

	return &cli.Command{
		Name:                   "bb",
		Usage:                  "Bitbucket Cloud CLI",
		Version:                Version,
		EnableShellCompletion:  true,
		Before: func(ctx context.Context, cmd *cli.Command) (context.Context, error) {
			f, err := cmdutil.NewFactory()
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", output.Error.Render(err.Error()))
				return ctx, err
			}
			return cmdutil.WithFactory(ctx, f), nil
		},
		Commands: []*cli.Command{
			authCmd.NewCmdAuth(),
			repoCmd.NewCmdRepo(),
			prCmd.NewCmdPR(),
			pipelineCmd.NewCmdPipeline(),
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
}
