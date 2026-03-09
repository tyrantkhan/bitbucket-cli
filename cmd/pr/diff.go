package pr

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/tyrantkhan/bb/internal/cmdutil"
	"github.com/tyrantkhan/bb/internal/output"
	"github.com/urfave/cli/v3"
)

func newCmdDiff() *cli.Command {
	return &cli.Command{
		Name:      "diff",
		Usage:     "Show the diff of a pull request",
		ArgsUsage: "<id>",
		Flags: []cli.Flag{
			cmdutil.WorkspaceFlag,
			cmdutil.RepoFlag,
			cmdutil.FormatFlag,
			&cli.BoolFlag{
				Name:  "stat",
				Usage: "Show only file-level summary",
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

			path := fmt.Sprintf("/2.0/repositories/%s/%s/pullrequests/%d/diff", workspace, repo, prID)

			resp, err := client.Get(path)
			if err != nil {
				return err
			}
			defer func() { _ = resp.Body.Close() }()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return fmt.Errorf("failed to read diff: %w", err)
			}

			diffText := string(body)

			format := cmdutil.GetFormat(ctx, cmd)
			if format == "json" {
				return output.RenderJSON(parseDiffFiles(diffText))
			}

			if cmd.Bool("stat") {
				printDiffStat(f, diffText)
				return nil
			}

			printColoredDiff(f, diffText)
			return nil
		},
	}
}

func printColoredDiff(f *cmdutil.Factory, diffText string) {
	scanner := bufio.NewScanner(strings.NewReader(diffText))
	for scanner.Scan() {
		line := scanner.Text()

		switch {
		case strings.HasPrefix(line, "diff "):
			fmt.Fprintln(f.IOOut, output.Bold.Render(line))
		case strings.HasPrefix(line, "@@"):
			fmt.Fprintln(f.IOOut, output.Cyan.Render(line))
		case strings.HasPrefix(line, "+"):
			fmt.Fprintln(f.IOOut, output.Green.Render(line))
		case strings.HasPrefix(line, "-"):
			fmt.Fprintln(f.IOOut, output.Red.Render(line))
		default:
			fmt.Fprintln(f.IOOut, line)
		}
	}
}

func printDiffStat(f *cmdutil.Factory, diffText string) {
	files := parseDiffFiles(diffText)

	if len(files) == 0 {
		fmt.Fprintln(f.IOOut, output.Muted.Render("No changes."))
		return
	}

	totalAdd := 0
	totalDel := 0

	for _, file := range files {
		adds := output.Green.Render(fmt.Sprintf("+%d", file.Additions))
		dels := output.Red.Render(fmt.Sprintf("-%d", file.Deletions))
		fmt.Fprintf(f.IOOut, " %s  %s %s\n", file.Path, adds, dels)
		totalAdd += file.Additions
		totalDel += file.Deletions
	}

	fmt.Fprintln(f.IOOut)
	fmt.Fprintf(f.IOOut, " %d file(s) changed, %s, %s\n",
		len(files),
		output.Green.Render(fmt.Sprintf("%d insertion(s)", totalAdd)),
		output.Red.Render(fmt.Sprintf("%d deletion(s)", totalDel)),
	)
}

type diffFile struct {
	Path      string `json:"path"`
	Additions int    `json:"additions"`
	Deletions int    `json:"deletions"`
	Patch     string `json:"patch"`
}

func parseDiffFiles(diffText string) []diffFile {
	var files []diffFile
	var current *diffFile
	var patchLines []string

	flushPatch := func() {
		if current != nil && len(patchLines) > 0 {
			current.Patch = strings.Join(patchLines, "\n")
			patchLines = nil
		}
	}

	scanner := bufio.NewScanner(strings.NewReader(diffText))
	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "diff --git") {
			flushPatch()
			parts := strings.Fields(line)
			name := ""
			if len(parts) >= 4 {
				name = strings.TrimPrefix(parts[3], "b/")
			}
			files = append(files, diffFile{Path: name})
			current = &files[len(files)-1]
			patchLines = append(patchLines, line)
			continue
		}

		if current == nil {
			continue
		}

		patchLines = append(patchLines, line)

		if strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++") {
			current.Additions++
		} else if strings.HasPrefix(line, "-") && !strings.HasPrefix(line, "---") {
			current.Deletions++
		}
	}
	flushPatch()

	if len(files) == 0 {
		return []diffFile{}
	}
	return files
}
