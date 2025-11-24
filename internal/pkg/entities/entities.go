package entities

import "time"

// Paper represents a paper in the arXiv dataset
type Paper struct {
	// ID of the paper (e.g., http://arxiv.org/abs/2511.17464v1)
	ID string `json:"id"`

	// Title of the paper
	Title string `json:"title"`

	// Summary/Abstract of the paper
	Summary string `json:"summary"`

	// Authors of the paper
	Authors []Author `json:"author"`

	// Publish date of the paper
	PublishDate time.Time `json:"publish_date"`

	// Updated date of the paper
	UpdatedDate time.Time `json:"updated_date"`

	// Links associated with the paper
	Links []Link `json:"links"`

	// Categories of the paper
	Categories []string `json:"categories"`
}

// Author represents an author of a paper
type Author struct {
	// Name of the author
	Name string `json:"name,omitempty"`

	// Affiliation of the author
	Affiliation string `json:"affiliation,omitempty"`

	// Country of the author
	Country string `json:"country,omitempty"`
}

// Link represents a link associated with a paper
type Link struct {
	Href string `json:"href"`
	Rel  string `json:"rel,omitempty"`
	Type string `json:"type,omitempty"`
}
