package api

// PaginatedResponse represents a Bitbucket paginated API response.
type PaginatedResponse[T any] struct {
	Size     int    `json:"size"`
	Page     int    `json:"page"`
	PageLen  int    `json:"pagelen"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Values   []T    `json:"values"`
}
