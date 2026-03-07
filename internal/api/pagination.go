package api

import "encoding/json"

// Paginate fetches all pages of a paginated endpoint up to the given limit.
// If limit <= 0, all pages are fetched.
func Paginate[T any](client *Client, path string, limit int) ([]T, error) {
	var all []T

	resp, err := client.Get(path) //nolint:bodyclose // closed by DecodeJSON
	if err != nil {
		return nil, err
	}

	var page PaginatedResponse[T]
	if err := DecodeJSON(resp, &page); err != nil {
		return nil, err
	}
	all = append(all, page.Values...)

	for page.Next != "" && (limit <= 0 || len(all) < limit) {
		resp, err := client.GetURL(page.Next) //nolint:bodyclose // closed by DecodeJSON
		if err != nil {
			return nil, err
		}

		page = PaginatedResponse[T]{}
		if err := DecodeJSON(resp, &page); err != nil {
			return nil, err
		}
		all = append(all, page.Values...)
	}

	if limit > 0 && len(all) > limit {
		all = all[:limit]
	}

	return all, nil
}

// PaginateRaw fetches all pages and returns raw JSON messages for custom decoding.
func PaginateRaw(client *Client, path string, limit int) ([]json.RawMessage, error) {
	return Paginate[json.RawMessage](client, path, limit)
}
