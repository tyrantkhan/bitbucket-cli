package pr

import "github.com/tyrantkhan/bb/internal/models"

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
