package auth

// Credentials holds the authentication credentials.
type Credentials struct {
	Username    string `json:"username"`
	AppPassword string `json:"app_password"`
	// Future OAuth fields
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

// Store is the interface for credential storage.
type Store interface {
	GetCredentials() (*Credentials, error)
	SetCredentials(creds *Credentials) error
	DeleteCredentials() error
}
