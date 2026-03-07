package auth

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
)

// Default OAuth consumer credentials. These are embedded in the CLI binary
// and are safe for public clients (same pattern as GitHub CLI).
// Override with BB_CLIENT_ID / BB_CLIENT_SECRET env vars or --client-id / --client-secret flags.
var (
	DefaultClientID     = ""
	DefaultClientSecret = ""
)

const (
	authorizeURL = "https://bitbucket.org/site/oauth2/authorize"
	tokenURL     = "https://bitbucket.org/site/oauth2/access_token"
)

// TokenResponse holds the response from the OAuth token endpoint.
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	Scopes       string `json:"scopes"`
}

// AuthorizeURL builds the OAuth authorization URL.
func AuthorizeURL(clientID string, port int) string {
	params := url.Values{
		"client_id":     {clientID},
		"response_type": {"code"},
		"redirect_uri":  {fmt.Sprintf("http://localhost:%d/callback", port)},
	}
	return authorizeURL + "?" + params.Encode()
}

// StartCallbackServer starts an HTTP server on a random port to receive the OAuth callback.
// It returns the port, a channel that will receive the authorization code, and any error.
func StartCallbackServer() (int, chan string, chan error, error) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, nil, nil, fmt.Errorf("failed to start callback server: %w", err)
	}

	port := listener.Addr().(*net.TCPAddr).Port
	codeCh := make(chan string, 1)
	errCh := make(chan error, 1)

	mux := http.NewServeMux()
	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code == "" {
			errMsg := r.URL.Query().Get("error_description")
			if errMsg == "" {
				errMsg = r.URL.Query().Get("error")
			}
			if errMsg == "" {
				errMsg = "no authorization code received"
			}
			w.Header().Set("Content-Type", "text/html")
			fmt.Fprintf(w, `<html><body><h2>Authentication Failed</h2><p>%s</p><p>You may close this tab.</p></body></html>`, errMsg)
			errCh <- fmt.Errorf("OAuth callback error: %s", errMsg)
			return
		}

		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, `<html><body><h2>Authentication Successful!</h2><p>You may close this tab and return to the terminal.</p></body></html>`)
		codeCh <- code
	})

	server := &http.Server{Handler: mux}
	go func() {
		if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	return port, codeCh, errCh, nil
}

// ExchangeCode exchanges an authorization code for tokens.
func ExchangeCode(clientID, clientSecret, code string, port int) (*TokenResponse, error) {
	data := url.Values{
		"grant_type":   {"authorization_code"},
		"code":         {code},
		"redirect_uri": {fmt.Sprintf("http://localhost:%d/callback", port)},
	}

	return tokenRequest(clientID, clientSecret, data)
}

// RefreshAccessToken refreshes an expired access token.
func RefreshAccessToken(clientID, clientSecret, refreshToken string) (*TokenResponse, error) {
	data := url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {refreshToken},
	}

	return tokenRequest(clientID, clientSecret, data)
}

func tokenRequest(clientID, clientSecret string, data url.Values) (*TokenResponse, error) {
	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(clientID, clientSecret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("token request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp struct {
			Error       string `json:"error"`
			Description string `json:"error_description"`
		}
		json.NewDecoder(resp.Body).Decode(&errResp)
		msg := errResp.Description
		if msg == "" {
			msg = errResp.Error
		}
		if msg == "" {
			msg = resp.Status
		}
		return nil, fmt.Errorf("token exchange failed: %s", msg)
	}

	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("failed to decode token response: %w", err)
	}

	return &tokenResp, nil
}
