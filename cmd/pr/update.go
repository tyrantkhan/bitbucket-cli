package pr

import (
	"fmt"
	"io"

	"github.com/tyrantkhan/bb/internal/api"
	"github.com/tyrantkhan/bb/internal/models"
	"github.com/tyrantkhan/bb/internal/output"
)

// buildPRUpdateBody builds a PUT body from an existing PR, preserving all
// fields that the Bitbucket API would otherwise drop.
func buildPRUpdateBody(pr models.PullRequest) map[string]interface{} {
	body := map[string]interface{}{
		"title":               pr.Title,
		"description":         pr.Description,
		"draft":               pr.Draft,
		"close_source_branch": pr.CloseSourceBranch,
		"source": map[string]interface{}{
			"branch": map[string]string{
				"name": pr.Source.Branch.Name,
			},
		},
		"destination": map[string]interface{}{
			"branch": map[string]string{
				"name": pr.Destination.Branch.Name,
			},
		},
	}

	if len(pr.Reviewers) > 0 {
		reviewers := make([]map[string]string, len(pr.Reviewers))
		for i, r := range pr.Reviewers {
			reviewers[i] = map[string]string{"uuid": r.UUID}
		}
		body["reviewers"] = reviewers
	}

	return body
}

// setDraftStatus fetches the PR, checks its current draft state, and updates
// it via PUT if needed. Used by both `bb pr ready` and `bb pr draft`.
func setDraftStatus(client *api.Client, out io.Writer, path string, prID int, draft bool) error {
	resp, err := client.Get(path)
	if err != nil {
		return err
	}

	var pr models.PullRequest
	if err := api.DecodeJSON(resp, &pr); err != nil {
		return fmt.Errorf("failed to decode pull request: %w", err)
	}

	if pr.Draft == draft {
		msg := fmt.Sprintf("Pull request #%d is already ready for review.", prID)
		if draft {
			msg = fmt.Sprintf("Pull request #%d is already a draft.", prID)
		}
		fmt.Fprintln(out, output.Muted.Render(msg))
		return nil
	}

	body := buildPRUpdateBody(pr)
	body["draft"] = draft

	resp, err = client.Put(path, body)
	if err != nil {
		return err
	}
	_ = resp.Body.Close()

	msg := fmt.Sprintf("Pull request #%d is now ready for review.", prID)
	if draft {
		msg = fmt.Sprintf("Pull request #%d is now a draft.", prID)
	}
	fmt.Fprintln(out, output.Success.Render(msg))

	return nil
}
