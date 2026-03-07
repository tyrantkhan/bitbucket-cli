package pr

import (
	"context"
	"fmt"
	"net/url"

	"github.com/tyrantkhan/bb/internal/api"
	"github.com/tyrantkhan/bb/internal/cmdutil"
	"github.com/tyrantkhan/bb/internal/git"
	"github.com/tyrantkhan/bb/internal/models"
	"github.com/tyrantkhan/bb/internal/output"
	"github.com/urfave/cli/v3"
)

type prStatusResult struct {
	CurrentBranch    []models.PullRequest `json:"current_branch"`
	CreatedByYou     []models.PullRequest `json:"created_by_you"`
	ReviewRequesting []models.PullRequest `json:"review_requesting"`
}

func newCmdStatus() *cli.Command {
	return &cli.Command{
		Name:  "status",
		Usage: "Show status of relevant pull requests",
		Flags: []cli.Flag{
			cmdutil.WorkspaceFlag,
			cmdutil.RepoFlag,
			cmdutil.FormatFlag,
		},
		Action: cmdutil.NoArgs(func(ctx context.Context, cmd *cli.Command) error {
			f := cmdutil.GetFactory(ctx)
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			workspace, repo, err := cmdutil.ResolveWorkspaceAndRepo(ctx, cmd)
			if err != nil {
				return err
			}

			// Get current user UUID.
			resp, err := client.Get("/2.0/user")
			if err != nil {
				return fmt.Errorf("failed to get current user: %w", err)
			}
			var user models.User
			if err := api.DecodeJSON(resp, &user); err != nil {
				return fmt.Errorf("failed to decode user: %w", err)
			}

			basePath := fmt.Sprintf("/2.0/repositories/%s/%s/pullrequests", workspace, repo)

			// Current branch PRs (soft-fail if not in a Bitbucket repo).
			detectedRepo, _ := git.DetectRepo()
			branch, branchErr := git.CurrentBranch()
			var branchPRs []models.PullRequest
			if branchErr == nil && branch != "" && detectedRepo != nil {
				q := fmt.Sprintf(`source.branch.name="%s"`, branch)
				path := fmt.Sprintf("%s?state=OPEN&q=%s", basePath, url.QueryEscape(q))
				branchPRs, err = api.Paginate[models.PullRequest](client, path, 0)
				if err != nil {
					return err
				}
			}

			// PRs created by the current user.
			q := fmt.Sprintf(`author.uuid="%s"`, user.UUID)
			path := fmt.Sprintf("%s?state=OPEN&q=%s", basePath, url.QueryEscape(q))
			createdPRs, err := api.Paginate[models.PullRequest](client, path, 0)
			if err != nil {
				return err
			}

			// PRs requesting review from the current user.
			q = fmt.Sprintf(`reviewers.uuid="%s"`, user.UUID)
			path = fmt.Sprintf("%s?state=OPEN&q=%s", basePath, url.QueryEscape(q))
			reviewPRs, err := api.Paginate[models.PullRequest](client, path, 0)
			if err != nil {
				return err
			}

			format := cmdutil.GetFormat(ctx, cmd)

			if format == "json" {
				return output.RenderJSON(prStatusResult{
					CurrentBranch:    branchPRs,
					CreatedByYou:     createdPRs,
					ReviewRequesting: reviewPRs,
				})
			}

			// Table output: three sections.
			fmt.Fprintln(f.IOOut)
			renderStatusSection(f, "Current branch", branch, branchErr, detectedRepo != nil, branchPRs)
			renderStatusSection(f, "Created by you", "", nil, true, createdPRs)
			renderStatusSection(f, "Requesting a code review from you", "", nil, true, reviewPRs)

			return nil
		}),
	}
}

func renderStatusSection(f *cmdutil.Factory, title string, branch string, branchErr error, isBitbucketRepo bool, prs []models.PullRequest) {
	fmt.Fprintln(f.IOOut, output.Header.Render(title))

	if title == "Current branch" {
		if branchErr != nil || branch == "" {
			fmt.Fprintln(f.IOOut, output.Muted.Render("  Not in a git repository"))
			fmt.Fprintln(f.IOOut)
			return
		}
		if !isBitbucketRepo {
			fmt.Fprintln(f.IOOut, output.Muted.Render("  Not in a Bitbucket repository"))
			fmt.Fprintln(f.IOOut)
			return
		}
		if len(prs) == 0 {
			fmt.Fprintf(f.IOOut, "  %s\n", output.Muted.Render(fmt.Sprintf("No pull request for branch %q", branch)))
			fmt.Fprintln(f.IOOut)
			return
		}
	}

	if len(prs) == 0 {
		fmt.Fprintln(f.IOOut, output.Muted.Render("  No pull requests"))
		fmt.Fprintln(f.IOOut)
		return
	}

	for _, pr := range prs {
		branchInfo := fmt.Sprintf("[%s → %s]", pr.Source.Branch.Name, pr.Destination.Branch.Name)
		summary := approvalSummary(pr)
		fmt.Fprintf(f.IOOut, "  #%-4d %s %s  %s  %s\n",
			pr.ID,
			pr.Title,
			output.Cyan.Render(branchInfo),
			output.StatusColor(pr.State).Render(pr.State),
			summary,
		)
	}
	fmt.Fprintln(f.IOOut)
}

func approvalSummary(pr models.PullRequest) string {
	var reviewers, approved, changesRequested int
	for _, p := range pr.Participants {
		if p.Role != "REVIEWER" {
			continue
		}
		reviewers++
		if p.Approved {
			approved++
		} else if p.State == "changes_requested" {
			changesRequested++
		}
	}

	if reviewers == 0 {
		return output.Muted.Render("No reviewers")
	}
	if changesRequested > 0 {
		return output.Red.Render("Changes requested")
	}
	if approved == reviewers {
		return output.Green.Render(fmt.Sprintf("%d/%d approved", approved, reviewers))
	}
	if approved > 0 {
		return output.Yellow.Render(fmt.Sprintf("%d/%d approved", approved, reviewers))
	}
	return output.Yellow.Render("Review required")
}
