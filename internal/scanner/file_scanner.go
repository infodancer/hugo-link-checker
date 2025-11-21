package scanner

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// EnumerateFiles recursively finds all files with the specified extension
// and returns a map of canonical paths to File structs to ensure uniqueness
func EnumerateFiles(rootDir string, extension string) (map[string]*File, error) {
	files := make(map[string]*File)
	
	// Normalize the extension to include the dot
	if !strings.HasPrefix(extension, ".") {
		extension = "." + extension
	}
	
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		// Skip directories
		if info.IsDir() {
			return nil
		}
		
		// Check if file has the desired extension
		if !strings.HasSuffix(strings.ToLower(path), strings.ToLower(extension)) {
			return nil
		}
		
		// Skip files beginning with a dot
		filename := filepath.Base(path)
		if strings.HasPrefix(filename, ".") {
			return nil
		}
		
		// Get canonical path to ensure uniqueness
		canonicalPath, err := filepath.Abs(path)
		if err != nil {
			return fmt.Errorf("failed to get canonical path for %s: %w", path, err)
		}
		
		// Clean the canonical path
		canonicalPath = filepath.Clean(canonicalPath)
		
		// Check if we've already seen this canonical path
		if _, exists := files[canonicalPath]; exists {
			// Skip duplicate files (e.g., symlinks pointing to same file)
			return nil
		}
		
		// Create new File struct
		file := &File{
			Path:          path,
			CanonicalPath: canonicalPath,
			Links:         make([]Link, 0),
		}
		
		files[canonicalPath] = file
		
		return nil
	})
	
	if err != nil {
		return nil, fmt.Errorf("failed to enumerate files: %w", err)
	}
	
	return files, nil
}

// GetFileList returns a slice of File pointers from the map for easier iteration
func GetFileList(fileMap map[string]*File) []*File {
	files := make([]*File, 0, len(fileMap))
	for _, file := range fileMap {
		files = append(files, file)
	}
	return files
}
