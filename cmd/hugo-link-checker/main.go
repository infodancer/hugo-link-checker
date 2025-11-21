package main

import (
    "flag"
    "fmt"
    "os"

    "github.com/infodancer/hugo-link-checker/internal/scanner"
    "github.com/infodancer/hugo-link-checker/internal/version"
)

func main() {
    var showVersion bool
    flag.BoolVar(&showVersion, "version", false, "Print version and exit")
    flag.Parse()

    if showVersion {
        fmt.Println("hugo-link-checker", version.Version)
        os.Exit(0)
    }

    // Example usage of the file scanner
    files, err := scanner.EnumerateFiles(".", ".md")
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error scanning files: %v\n", err)
        os.Exit(1)
    }
    
    fmt.Printf("Found %d unique markdown files\n", len(files))
    for _, file := range scanner.GetFileList(files) {
        fmt.Printf("File: %s (canonical: %s)\n", file.Path, file.CanonicalPath)
    }
}
