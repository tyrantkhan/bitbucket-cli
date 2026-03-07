package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"

	"github.com/tyrantkhan/bb/internal/auth"
)

const (
	baseURL   = "https://api.bitbucket.org"
	userAgent = "bb-cli/0.1.0"
)

// slugPattern validates workspace slugs, repo slugs, and similar path segments.
var slugPattern = regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)

// ValidateSlug checks that a path segment is safe for URL interpolation.
func ValidateSlug(name, value string) error {
	if !slugPattern.MatchString(value) {
		return fmt.Errorf("invalid %s: %q (must match [a-zA-Z0-9._-]+)", name, value)
	}
	return nil
}

// Client is the Bitbucket API client.
type Client struct {
	http        *http.Client
	username    string
	password    string
	authMethod  string
	bearerToken string
}

// NewClient creates a new API client with the given credentials.
func NewClient(creds *auth.Credentials) *Client {
	c := &Client{
		http: &http.Client{},
	}

	if creds.IsOAuth() {
		c.authMethod = "oauth"
		c.bearerToken = creds.AccessToken
	} else {
		c.authMethod = "api_token"
		c.username = creds.Username
		c.password = creds.APIToken
	}

	return c
}

// Get performs a GET request to the given API path.
func (c *Client) Get(path string) (*http.Response, error) {
	return c.do("GET", path, nil)
}

// GetURL performs a GET request to an absolute URL (for pagination).
func (c *Client) GetURL(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	c.setHeaders(req)
	return c.doRequest(req)
}

// Post performs a POST request with a JSON body.
func (c *Client) Post(path string, body interface{}) (*http.Response, error) {
	return c.do("POST", path, body)
}

// Put performs a PUT request with a JSON body.
func (c *Client) Put(path string, body interface{}) (*http.Response, error) {
	return c.do("PUT", path, body)
}

// Delete performs a DELETE request.
func (c *Client) Delete(path string) (*http.Response, error) {
	return c.do("DELETE", path, nil)
}

func (c *Client) do(method, path string, body interface{}) (*http.Response, error) {
	url := baseURL + path

	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewReader(data)
	}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, err
	}

	c.setHeaders(req)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return c.doRequest(req)
}

func (c *Client) setHeaders(req *http.Request) {
	if c.authMethod == "oauth" {
		req.Header.Set("Authorization", "Bearer "+c.bearerToken)
	} else {
		req.SetBasicAuth(c.username, c.password)
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "application/json")
}

func (c *Client) doRequest(req *http.Request) (*http.Response, error) {
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		defer func() { _ = resp.Body.Close() }()
		return nil, parseErrorResponse(resp)
	}

	return resp, nil
}

// DecodeJSON reads a JSON response body into the target struct.
func DecodeJSON(resp *http.Response, target interface{}) error {
	defer func() { _ = resp.Body.Close() }()
	return json.NewDecoder(resp.Body).Decode(target)
}
