package models

import "time"

// User represents a Bitbucket user.
type User struct {
	DisplayName string `json:"display_name"`
	UUID        string `json:"uuid"`
	Nickname    string `json:"nickname"`
	AccountID   string `json:"account_id"`
	Links       Links  `json:"links"`
}

// Workspace represents a Bitbucket workspace.
type Workspace struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// Project represents a Bitbucket project.
type Project struct {
	Key  string `json:"key"`
	Name string `json:"name"`
	UUID string `json:"uuid"`
}

// Branch represents a branch reference.
type Branch struct {
	Name string `json:"name"`
}

// Commit represents a commit reference.
type Commit struct {
	Hash    string `json:"hash"`
	Message string `json:"message"`
	Date    string `json:"date"`
	Author  *User  `json:"author"`
}

// Links holds common link structures.
type Links struct {
	Self   *Link   `json:"self,omitempty"`
	HTML   *Link   `json:"html,omitempty"`
	Clone  []Clone `json:"clone,omitempty"`
	Avatar *Link   `json:"avatar,omitempty"`
}

// Link is a single link.
type Link struct {
	Href string `json:"href"`
}

// Clone represents a clone link with its protocol name.
type Clone struct {
	Href string `json:"href"`
	Name string `json:"name"`
}

// Participant represents a PR participant.
type Participant struct {
	User     User   `json:"user"`
	Role     string `json:"role"` // "REVIEWER", "PARTICIPANT"
	Approved bool   `json:"approved"`
	State    string `json:"state"` // "approved", "changes_requested", etc.
}

// FormatTime formats a Bitbucket timestamp for display.
func FormatTime(t string) string {
	parsed, err := time.Parse(time.RFC3339Nano, t)
	if err != nil {
		parsed, err = time.Parse("2006-01-02T15:04:05.000000+00:00", t)
		if err != nil {
			return t
		}
	}
	return parsed.Local().Format("2006-01-02 15:04")
}
