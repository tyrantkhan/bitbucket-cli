package models

// SearchCodeResult represents a single code search result from the Bitbucket API.
type SearchCodeResult struct {
	Type              string               `json:"type"`
	ContentMatchCount int                  `json:"content_match_count"`
	ContentMatches    []SearchContentMatch `json:"content_matches"`
	PathMatches       []SearchSegment      `json:"path_matches"`
	File              SearchFile           `json:"file"`
}

// SearchContentMatch holds a set of matched lines within a file.
type SearchContentMatch struct {
	Lines []SearchLine `json:"lines"`
}

// SearchLine represents a single line in a content match.
type SearchLine struct {
	Line     int             `json:"line"`
	Segments []SearchSegment `json:"segments"`
}

// SearchSegment represents a segment of text, possibly a match.
type SearchSegment struct {
	Text  string `json:"text"`
	Match bool   `json:"match"`
}

// SearchFile describes the file that matched a code search.
type SearchFile struct {
	Type  string `json:"type"`
	Path  string `json:"path"`
	Links Links  `json:"links"`
}
