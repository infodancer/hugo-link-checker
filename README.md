# hugo-link-checker

Command-line tool to check links in Hugo-based websites.

This repository contains a Go-based CLI `hugo-link-checker` and CI workflow
to produce compiled binaries for Linux, macOS, and Windows using GitHub Actions.

Quick start (local build):

```bash
make build
./hugo-link-checker -version
```

Build cross-platform locally:

```bash
make build-all
ls dist
```

Next steps:
- Implement link crawling and checking logic in `internal/checker`.
- Add flags and options to `cmd/hugo-link-checker/main.go`.
- Wire up tests and CI checks.
