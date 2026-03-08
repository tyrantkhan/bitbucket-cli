package search

import (
	"context"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/tyrantkhan/bb/internal/api"
	"github.com/tyrantkhan/bb/internal/cmdutil"
	"github.com/tyrantkhan/bb/internal/models"
	"github.com/tyrantkhan/bb/internal/output"
	"github.com/urfave/cli/v3"
)

func newCmdCode() *cli.Command {
	return &cli.Command{
		Name:      "code",
		Usage:     "Search for code across repositories in a workspace",
		ArgsUsage: "<query>",
		Flags: []cli.Flag{
			cmdutil.WorkspaceFlag,
			cmdutil.FormatFlag,
			&cli.IntFlag{
				Name:  "limit",
				Usage: "Maximum number of results",
				Value: 20,
			},
			&cli.StringFlag{
				Name:    "repo",
				Aliases: []string{"R"},
				Usage:   "Filter by repository slug (appends repo:<val> to query)",
			},
			&cli.StringFlag{
				Name:    "extension",
				Aliases: []string{"e"},
				Usage:   "Filter by file extension (appends ext:<val> to query)",
			},
			&cli.StringFlag{
				Name:    "language",
				Aliases: []string{"l"},
				Usage:   "Filter by language (appends lang:<val> to query)",
			},
			&cli.StringFlag{
				Name:    "path",
				Aliases: []string{"p"},
				Usage:   "Filter by file path (appends path:<val> to query)",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if cmd.Args().Len() < 1 {
				_ = cli.ShowSubcommandHelp(cmd)
				return cmdutil.ErrShowedUsage
			}

			f := cmdutil.GetFactory(ctx)
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			workspace, err := cmdutil.ResolveWorkspace(ctx, cmd)
			if err != nil {
				return err
			}

			query := buildQuery(cmd)
			limit := int(cmd.Int("limit"))

			path := fmt.Sprintf("/2.0/workspaces/%s/search/code?search_query=%s",
				workspace, url.QueryEscape(query))

			results, err := api.Paginate[models.SearchCodeResult](client, path, limit)
			if err != nil {
				return err
			}

			format := cmdutil.GetFormat(ctx, cmd)
			if format == "json" {
				return output.RenderJSON(results)
			}

			renderCodeResults(results)
			return nil
		},
	}
}

// buildQuery constructs the search query from the positional arg and convenience flags.
func buildQuery(cmd *cli.Command) string {
	parts := []string{cmd.Args().First()}

	if v := cmd.String("repo"); v != "" {
		parts = append(parts, "repo:"+v)
	}
	if v := cmd.String("extension"); v != "" {
		parts = append(parts, "ext:"+v)
	}
	if v := cmd.String("language"); v != "" {
		parts = append(parts, "lang:"+v)
	}
	if v := cmd.String("path"); v != "" {
		parts = append(parts, "path:"+v)
	}

	return strings.Join(parts, " ")
}

// renderCodeResults prints code search results in a human-friendly format.
func renderCodeResults(results []models.SearchCodeResult) {
	if len(results) == 0 {
		fmt.Println(output.Muted.Render("No code results found."))
		return
	}

	fmt.Printf("Showing %d code result%s\n", len(results), plural(len(results)))

	for _, r := range results {
		repoSlug := extractRepoSlug(r.File)
		filePath := sanitize(r.File.Path)
		if repoSlug != "" {
			filePath = sanitize(repoSlug) + "/" + filePath
		}

		fmt.Println()
		fmt.Println(output.Bold.Render(filePath))

		for _, cm := range r.ContentMatches {
			for _, line := range cm.Lines {
				lineNum := output.Muted.Render(fmt.Sprintf("%4d: ", line.Line))
				var text strings.Builder
				for _, seg := range line.Segments {
					clean := sanitize(seg.Text)
					if seg.Match {
						text.WriteString(output.Bold.Render(clean))
					} else {
						text.WriteString(clean)
					}
				}
				fmt.Printf("  %s%s\n", lineNum, text.String())
			}
		}
	}
}

// extractRepoSlug extracts the repository slug from the file's self link.
// The link typically looks like: https://api.bitbucket.org/2.0/repositories/{workspace}/{repo}/src/...
func extractRepoSlug(f models.SearchFile) string {
	if f.Links.Self == nil {
		return ""
	}
	parts := strings.Split(f.Links.Self.Href, "/")
	for i, p := range parts {
		if p == "repositories" && i+2 < len(parts) {
			return parts[i+2]
		}
	}
	return ""
}

var ansiEscapeRe = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)

// sanitize strips ANSI escape sequences from untrusted text to prevent terminal injection.
func sanitize(s string) string {
	return ansiEscapeRe.ReplaceAllString(s, "")
}

func plural(n int) string {
	if n == 1 {
		return ""
	}
	return "s"
}
