package models

// PullRequest represents a Bitbucket pull request.
type PullRequest struct {
	ID                int           `json:"id"`
	Title             string        `json:"title"`
	Description       string        `json:"description"`
	State             string        `json:"state"` // OPEN, MERGED, DECLINED, SUPERSEDED
	CreatedOn         string        `json:"created_on"`
	UpdatedOn         string        `json:"updated_on"`
	Author            User          `json:"author"`
	Source            PREndpoint    `json:"source"`
	Destination       PREndpoint    `json:"destination"`
	CloseSourceBranch bool          `json:"close_source_branch"`
	MergeCommit       *Commit       `json:"merge_commit"`
	Reviewers         []User        `json:"reviewers"`
	Participants      []Participant `json:"participants"`
	TaskCount         int           `json:"task_count"`
	CommentCount      int           `json:"comment_count"`
	Reason            string        `json:"reason"`
	Links             Links         `json:"links"`
}

// PREndpoint represents a source or destination branch endpoint.
type PREndpoint struct {
	Branch     Branch      `json:"branch"`
	Commit     Commit      `json:"commit"`
	Repository *Repository `json:"repository"`
}
