package models

// Repository represents a Bitbucket repository.
type Repository struct {
	UUID        string    `json:"uuid"`
	Name        string    `json:"name"`
	FullName    string    `json:"full_name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	IsPrivate   bool      `json:"is_private"`
	Language    string    `json:"language"`
	CreatedOn   string    `json:"created_on"`
	UpdatedOn   string    `json:"updated_on"`
	Size        int64     `json:"size"`
	ForkPolicy  string    `json:"fork_policy"`
	SCM         string    `json:"scm"`
	Project     *Project  `json:"project"`
	Owner       *User     `json:"owner"`
	Workspace   Workspace `json:"workspace"`
	MainBranch  *Branch   `json:"mainbranch"`
	Links       Links     `json:"links"`
}

// Visibility returns "private" or "public".
func (r *Repository) Visibility() string {
	if r.IsPrivate {
		return "private"
	}
	return "public"
}

// CloneURL returns the clone URL for the given protocol ("ssh" or "https").
func (r *Repository) CloneURL(protocol string) string {
	for _, c := range r.Links.Clone {
		if c.Name == protocol {
			return c.Href
		}
	}
	return ""
}
