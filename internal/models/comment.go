package models

// Comment represents a pull request comment.
type Comment struct {
	ID      int  `json:"id"`
	Parent  *ID  `json:"parent,omitempty"`
	Deleted bool `json:"deleted"`
	Content struct {
		Raw    string `json:"raw"`
		Markup string `json:"markup"`
		HTML   string `json:"html"`
	} `json:"content"`
	Inline *InlineRef `json:"inline,omitempty"`
	Author User       `json:"user"`

	CreatedOn string `json:"created_on"`
	UpdatedOn string `json:"updated_on"`
	Links     Links  `json:"links"`
}

// InlineRef marks a comment as inline on a specific file/line.
type InlineRef struct {
	Path string `json:"path"`
	From *int   `json:"from,omitempty"` // old line
	To   *int   `json:"to,omitempty"`   // new line
}

// ID is used for parent reference.
type ID struct {
	ID int `json:"id"`
}
