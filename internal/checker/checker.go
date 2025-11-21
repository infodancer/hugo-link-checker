package checker

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/infodancer/hugo-link-checker/internal/scanner"
)

// CheckLinks validates all links in the provided files
func CheckLinks(files []*scanner.File, rootDir string) error {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	for _, file := range files {
		for i := range file.Links {
			link := &file.Links[i]
			
			if link.Type == scanner.LinkTypeExternal {
				err := checkExternalLink(client, link)
				if err != nil {
					return fmt.Errorf("error checking external link %s: %v", link.URL, err)
				}
			} else {
				err := checkInternalLink(link, rootDir)
				if err != nil {
					return fmt.Errorf("error checking internal link %s: %v", link.URL, err)
				}
			}
			
			link.LastChecked = time.Now()
		}
	}
	
	return nil
}

func checkExternalLink(client *http.Client, link *scanner.Link) error {
	resp, err := client.Head(link.URL)
	if err != nil {
		// Try GET if HEAD fails
		resp, err = client.Get(link.URL)
		if err != nil {
			link.StatusCode = 0
			link.ErrorMessage = err.Error()
			return nil
		}
	}
	defer resp.Body.Close()
	
	link.StatusCode = resp.StatusCode
	if resp.StatusCode >= 400 {
		link.ErrorMessage = fmt.Sprintf("HTTP %d", resp.StatusCode)
	} else {
		link.ErrorMessage = ""
	}
	
	return nil
}

func checkInternalLink(link *scanner.Link, rootDir string) error {
	// Clean and resolve the path
	linkPath := link.URL
	
	// Remove fragment identifier
	if idx := strings.Index(linkPath, "#"); idx != -1 {
		linkPath = linkPath[:idx]
	}
	
	// Remove query parameters
	if idx := strings.Index(linkPath, "?"); idx != -1 {
		linkPath = linkPath[:idx]
	}
	
	// Skip empty paths (fragment-only links)
	if linkPath == "" {
		link.StatusCode = 200
		link.ErrorMessage = ""
		return nil
	}
	
	// Resolve relative path
	var fullPath string
	if filepath.IsAbs(linkPath) {
		fullPath = filepath.Join(rootDir, linkPath)
	} else {
		fullPath = filepath.Join(rootDir, linkPath)
	}
	
	// Check if file exists
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		link.StatusCode = 404
		link.ErrorMessage = "File not found"
	} else if err != nil {
		link.StatusCode = 0
		link.ErrorMessage = err.Error()
	} else {
		link.StatusCode = 200
		link.ErrorMessage = ""
	}
	
	return nil
}

// CountBrokenLinks returns the number of broken links across all files
func CountBrokenLinks(files []*scanner.File) int {
	count := 0
	for _, file := range files {
		for _, link := range file.Links {
			if link.StatusCode >= 400 || (link.StatusCode == 0 && link.ErrorMessage != "") {
				count++
			}
		}
	}
	return count
}
