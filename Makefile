.PHONY: build build-all clean
BINARY := hugo-link-checker
DIST := dist

build:
	go build -o $(BINARY) ./cmd/hugo-link-checker

build-all:
	@echo "Building for multiple OS/ARCH..."
	@rm -rf $(DIST) && mkdir -p $(DIST)
	@for os_arch in "linux/amd64" "linux/arm64" "darwin/amd64" "darwin/arm64" "windows/amd64"; do \
	  os=$${os_arch%/*}; arch=$${os_arch#*/}; \
	  ext=""; if [ "$${os}" = "windows" ]; then ext=".exe"; fi; \
	  echo "Building $$os/$$arch"; \
	  CGO_ENABLED=0 GOOS=$$os GOARCH=$$arch go build -o $(DIST)/$(BINARY)-$$os-$$arch$$ext ./cmd/hugo-link-checker; \
	done

clean:
	rm -f $(BINARY)
	rm -rf $(DIST)
