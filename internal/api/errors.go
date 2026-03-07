package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// APIError represents an error from the Bitbucket API.
type APIError struct {
	StatusCode int
	Message    string
	Detail     string
}

func (e *APIError) Error() string {
	if e.Detail != "" {
		return fmt.Sprintf("API error %d: %s — %s", e.StatusCode, e.Message, e.Detail)
	}
	return fmt.Sprintf("API error %d: %s", e.StatusCode, e.Message)
}

// bitbucketErrorResponse represents the Bitbucket API error JSON structure.
type bitbucketErrorResponse struct {
	Error struct {
		Message string `json:"message"`
		Detail  string `json:"detail"`
	} `json:"error"`
}

// parseErrorResponse parses an HTTP response into an APIError.
func parseErrorResponse(resp *http.Response) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &APIError{
			StatusCode: resp.StatusCode,
			Message:    friendlyMessage(resp.StatusCode),
		}
	}

	var bbErr bitbucketErrorResponse
	if json.Unmarshal(body, &bbErr) == nil && bbErr.Error.Message != "" {
		return &APIError{
			StatusCode: resp.StatusCode,
			Message:    bbErr.Error.Message,
			Detail:     bbErr.Error.Detail,
		}
	}

	return &APIError{
		StatusCode: resp.StatusCode,
		Message:    friendlyMessage(resp.StatusCode),
	}
}

func friendlyMessage(code int) string {
	switch code {
	case 401:
		return "authentication failed. Check your credentials or run 'bb auth login'"
	case 403:
		return "permission denied. Your app password may lack the required scope"
	case 404:
		return "resource not found. Check the workspace, repository, or ID"
	case 429:
		return "rate limited. Please wait and try again"
	default:
		return http.StatusText(code)
	}
}
