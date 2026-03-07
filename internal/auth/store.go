package auth

// Credentials holds the authentication credentials.
type Credentials struct {
	Username    string `json:"username"`
	APIToken string `json:"api_token,omitempty"`

	// OAuth fields
	AuthMethod   string `json:"auth_method,omitempty"` // "api_token" or "oauth"
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	ExpiresAt    int64  `json:"expires_at,omitempty"`
	ClientID     string `json:"client_id,omitempty"`
	ClientSecret string `json:"client_secret,omitempty"`
}

// IsOAuth returns true if the credentials use OAuth authentication.
func (c *Credentials) IsOAuth() bool {
	return c.AuthMethod == "oauth"
}

// Store is the interface for credential storage.
type Store interface {
	GetCredentials() (*Credentials, error)
	SetCredentials(creds *Credentials) error
	DeleteCredentials() error
}
