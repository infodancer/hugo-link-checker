package scanner

import (
	"net/url"
	"path/filepath"
	"strings"
	"time"
)

// LinkType represents whether a link is internal or external
type LinkType int

const (
	LinkTypeInternal LinkType = iota
	LinkTypeExternal
)

// Link represents a link found in a file
type Link struct {
	URL          string    `json:"url"`
	Type         LinkType  `json:"type"`
	LastChecked  time.Time `json:"last_checked"`
	StatusCode   int       `json:"status_code"`
	ErrorMessage string    `json:"error_message,omitempty"`
}

// File represents a file and its links
type File struct {
	Path          string `json:"path"`
	CanonicalPath string `json:"canonical_path"`
	Links         []Link `json:"links"`
}

// isInternalLink determines if a link is internal (relative) or external
func isInternalLink(linkURL string) bool {
	// Parse the URL
	u, err := url.Parse(linkURL)
	if err != nil {
		// If we can't parse it, treat as internal for safety
		return true
	}
	
	// If it has a scheme (http, https, etc.) or host, it's external
	if u.Scheme != "" || u.Host != "" {
		return false
	}
	
	// Otherwise it's a relative/internal link
	return true
}

// NewLink creates a new Link with the appropriate type
func NewLink(linkURL string) Link {
	linkType := LinkTypeInternal
	if !isInternalLink(linkURL) {
		linkType = LinkTypeExternal
	}
	
	return Link{
		URL:  linkURL,
		Type: linkType,
	}
}
