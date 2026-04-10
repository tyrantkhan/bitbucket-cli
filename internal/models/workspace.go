package models

// WorkspaceMembership represents a member's membership in a workspace.
type WorkspaceMembership struct {
	User      User      `json:"user"`
	Workspace Workspace `json:"workspace"`
}
