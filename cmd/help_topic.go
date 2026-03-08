package cmd

import (
	"fmt"
	"strings"

	"github.com/urfave/cli/v3"
)

type helpTopic struct {
	name  string
	short string
	long  string
}

var helpTopics = []helpTopic{
	{
		name:  "environment",
		short: "Environment variables used by bb",
		long: `bb respects the following environment variables:

  XDG_CONFIG_HOME     Base directory for configuration files. When set, bb
                      stores credentials and config under
                      $XDG_CONFIG_HOME/bb/ instead of ~/.config/bb/.

  BB_CLIENT_ID        OAuth2 client ID. When set alongside BB_CLIENT_SECRET,
                      bb uses these credentials instead of the built-in
                      OAuth consumer.

  BB_CLIENT_SECRET    OAuth2 client secret. Used together with BB_CLIENT_ID.

  NO_COLOR            When set to any value, disables colour output.`,
	},
	{
		name:  "formatting",
		short: "Output formatting options",
		long: `Many bb commands support different output formats controlled by the
--format flag:

  --format table      Default. Human-readable table output with colours.
  --format json       Machine-readable JSON output for scripting.

Some commands also support:

  --web               Open the resource in your default web browser instead
                      of printing output.

The default format can be changed in ~/.config/bb/config.yml:

  format: json`,
	},
	{
		name:  "oauth",
		short: "Using a custom OAuth consumer",
		long: `By default, bb authenticates using a built-in OAuth consumer. If you
prefer to use your own (for security, compliance, or organizational
reasons), you can create a custom OAuth consumer in Bitbucket and
configure bb to use it.

CREATING AN OAUTH CONSUMER

  1. Go to your Bitbucket workspace settings:
     https://bitbucket.org/<workspace>/workspace/settings/api

  2. Click "Add consumer".

  3. Fill in:
     - Name:         anything (e.g. "bb CLI")
     - Callback URL: http://localhost/callback
     - Permissions:  Account (Read), Repositories (Read/Write),
                     Pull Requests (Read/Write), Pipelines (Read/Write)

  4. Save and note the Key (client ID) and Secret (client secret).

CONFIGURING BB

  Option 1 — environment variables (recommended):

    export BB_CLIENT_ID="your-key"
    export BB_CLIENT_SECRET="your-secret"
    bb auth login --web

  Option 2 — command-line flags:

    bb auth login --web --client-id "your-key" --client-secret "your-secret"

  The client ID and secret are stored alongside your tokens in
  ~/.config/bb/credentials.json so that token refresh works automatically.

NOTES

  - The callback URL in your consumer must start with http://localhost.
    bb starts a local server on a random port to receive the OAuth code.

  - OAuth consumers are scoped to a workspace. If you work across multiple
    workspaces, the consumer must be created in each one — or use the
    built-in consumer which works globally.

  - If you switch from a custom consumer back to the built-in one, run
    bb auth logout first to clear the stored client credentials.`,
	},
	{
		name:  "exit-codes",
		short: "Exit codes used by bb",
		long: `bb uses the following exit codes:

  0   Success.
  1   An error occurred (API failure, invalid arguments, etc.).`,
	},
	{
		name:  "reference",
		short: "Auto-generated command reference",
		// long is generated dynamically from the command tree
	},
}

func newHelpTopicCommands(root *cli.Command) []*cli.Command {
	cmds := make([]*cli.Command, 0, len(helpTopics))
	for _, t := range helpTopics {
		long := t.long
		if t.name == "reference" {
			long = referenceText(root)
		}
		cmds = append(cmds, &cli.Command{
			Name:        t.name,
			Usage:       t.short,
			Description: long,
			Hidden:      true,
		})
	}
	return cmds
}

func referenceText(root *cli.Command) string {
	var b strings.Builder
	b.WriteString("bb command reference:\n")
	walkCommands(&b, root.Commands, "bb", 0)
	return b.String()
}

func walkCommands(b *strings.Builder, cmds []*cli.Command, prefix string, depth int) {
	for _, c := range cmds {
		if c.Hidden {
			continue
		}
		indent := strings.Repeat("  ", depth)
		full := fmt.Sprintf("%s %s", prefix, c.Name)

		// Command name with optional args usage.
		header := full
		if c.ArgsUsage != "" {
			header += " " + c.ArgsUsage
		}
		fmt.Fprintf(b, "\n%s%s\n", indent, header)

		if c.Usage != "" {
			fmt.Fprintf(b, "%s  %s\n", indent, c.Usage)
		}

		// Aliases.
		if len(c.Aliases) > 0 {
			fmt.Fprintf(b, "%s  Aliases: %s\n", indent, strings.Join(c.Aliases, ", "))
		}

		// Flags (skip help which is always present).
		for _, fl := range c.Flags {
			names := fl.Names()
			if len(names) == 0 {
				continue
			}
			if names[0] == "help" {
				continue
			}

			var flagStr string
			short := ""
			long := names[0]
			for _, n := range names {
				if len(n) == 1 {
					short = n
				} else {
					long = n
				}
			}
			if short != "" {
				flagStr = fmt.Sprintf("-%s, --%s", short, long)
			} else {
				flagStr = fmt.Sprintf("    --%s", long)
			}

			usage := ""
			if du, ok := fl.(interface{ GetUsage() string }); ok {
				usage = du.GetUsage()
			}

			if usage != "" {
				fmt.Fprintf(b, "%s    %-24s %s\n", indent, flagStr, usage)
			} else {
				fmt.Fprintf(b, "%s    %s\n", indent, flagStr)
			}
		}

		if len(c.Commands) > 0 {
			walkCommands(b, c.Commands, full, depth+1)
		}
	}
}

func helpTopicsSection() string {
	var b strings.Builder
	for _, t := range helpTopics {
		fmt.Fprintf(&b, "   %-16s %s\n", t.name, t.short)
	}
	return b.String()
}
