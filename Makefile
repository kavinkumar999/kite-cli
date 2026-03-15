# Kite CLI Makefile

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME ?= $(shell date -u '+%Y-%m-%d %H:%M:%S')
GIT_COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

LDFLAGS := -s -w \
	-X 'github.com/kavinkumar999/kite-cli/cmd.Version=$(VERSION)' \
	-X 'github.com/kavinkumar999/kite-cli/cmd.BuildTime=$(BUILD_TIME)' \
	-X 'github.com/kavinkumar999/kite-cli/cmd.GitCommit=$(GIT_COMMIT)'

.PHONY: build install clean test release

# Build binary
build:
	go build -ldflags="$(LDFLAGS)" -o kite .

# Install to ~/bin
install: build
	mkdir -p ~/bin
	cp kite ~/bin/kite
	@echo "Installed to ~/bin/kite"

# Clean build artifacts
clean:
	rm -f kite
	rm -rf dist/

# Run tests
test:
	go test -v ./...

# Build for release (multiple platforms)
release: clean
	@mkdir -p dist
	@echo "Building for darwin/amd64..."
	GOOS=darwin GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o dist/kite_darwin_amd64 .
	@echo "Building for darwin/arm64..."
	GOOS=darwin GOARCH=arm64 go build -ldflags="$(LDFLAGS)" -o dist/kite_darwin_arm64 .
	@echo "Building for linux/amd64..."
	GOOS=linux GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o dist/kite_linux_amd64 .
	@echo "Building for linux/arm64..."
	GOOS=linux GOARCH=arm64 go build -ldflags="$(LDFLAGS)" -o dist/kite_linux_arm64 .
	@echo "Release binaries built in dist/"

# Show version info
version:
	@echo "Version: $(VERSION)"
	@echo "Build Time: $(BUILD_TIME)"
	@echo "Git Commit: $(GIT_COMMIT)"
