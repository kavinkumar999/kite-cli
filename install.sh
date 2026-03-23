#!/bin/bash
set -e

# Kite CLI Installer
# Usage: curl -sSL https://raw.githubusercontent.com/kavin/kite-cli/main/install.sh | bash

REPO="kavinkumar999/kite-cli"
INSTALL_DIR="$HOME/bin"
BINARY_NAME="kite"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

info() { echo -e "${GREEN}[INFO]${NC} $1"; }
warn() { echo -e "${YELLOW}[WARN]${NC} $1"; }
error() { echo -e "${RED}[ERROR]${NC} $1"; cleanup; exit 1; }

# Cleanup temp files on exit
cleanup() {
    rm -f /tmp/kite /tmp/kite.tar.gz 2>/dev/null || true
    rm -rf /tmp/kite-cli-* 2>/dev/null || true
}

# Set trap to cleanup on script exit
trap cleanup EXIT

# Detect OS and architecture
detect_platform() {
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    ARCH=$(uname -m)

    case $ARCH in
        x86_64)  ARCH="amd64" ;;
        aarch64) ARCH="arm64" ;;
        arm64)   ARCH="arm64" ;;
        *)       error "Unsupported architecture: $ARCH" ;;
    esac

    case $OS in
        darwin) OS="darwin" ;;
        linux)  OS="linux" ;;
        *)      error "Unsupported OS: $OS" ;;
    esac

    PLATFORM="${OS}_${ARCH}"
    info "Detected platform: $PLATFORM"
}

# Check if Go is installed
check_go() {
    if command -v go &> /dev/null; then
        GO_VERSION=$(go version | awk '{print $3}')
        info "Found Go: $GO_VERSION"
        return 0
    fi
    return 1
}

# Try to download pre-built binary
install_binary() {
    LATEST_RELEASE=$(curl -sSL "https://api.github.com/repos/${REPO}/releases/latest" 2>/dev/null | grep '"tag_name"' | sed -E 's/.*"([^"]+)".*/\1/' || echo "")
    
    if [ -z "$LATEST_RELEASE" ]; then
        warn "No releases found. Building from source..."
        return 1
    fi

    DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${LATEST_RELEASE}/kite_${PLATFORM}.tar.gz"
    
    info "Downloading ${BINARY_NAME} ${LATEST_RELEASE}..."
    
    # Use -f to fail on HTTP errors (404, etc.)
    if ! curl -fsSL "$DOWNLOAD_URL" -o /tmp/kite.tar.gz 2>/dev/null; then
        warn "No pre-built binary for ${PLATFORM}. Building from source..."
        rm -f /tmp/kite.tar.gz 2>/dev/null
        return 1
    fi
    
    # Verify it's a valid gzip file
    if ! gzip -t /tmp/kite.tar.gz 2>/dev/null; then
        warn "Downloaded file is not valid. Building from source..."
        rm -f /tmp/kite.tar.gz 2>/dev/null
        return 1
    fi
    
    tar -xzf /tmp/kite.tar.gz -C /tmp || return 1
    rm -f /tmp/kite.tar.gz
    
    # Verify binary exists after extraction
    if [ ! -f /tmp/kite ]; then
        warn "Binary not found in archive. Building from source..."
        return 1
    fi
    
    return 0
}

# Build from source
build_from_source() {
    if ! check_go; then
        error "Go is required to build from source. Install Go from https://golang.org/dl/"
    fi

    info "Building from source..."
    
    TEMP_DIR=$(mktemp -d)
    cd "$TEMP_DIR"
    
    info "Cloning repository..."
    git clone --depth 1 "https://github.com/${REPO}.git" kite-cli
    cd kite-cli
    
    info "Building..."
    VERSION=$(git describe --tags --always 2>/dev/null || echo "dev")
    BUILD_TIME=$(date -u '+%Y-%m-%d %H:%M:%S')
    GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
    
    LDFLAGS="-s -w"
    LDFLAGS="$LDFLAGS -X 'github.com/kavinkumar999/kite-cli/cmd.Version=$VERSION'"
    LDFLAGS="$LDFLAGS -X 'github.com/kavinkumar999/kite-cli/cmd.BuildTime=$BUILD_TIME'"
    LDFLAGS="$LDFLAGS -X 'github.com/kavinkumar999/kite-cli/cmd.GitCommit=$GIT_COMMIT'"
    
    go build -ldflags="$LDFLAGS" -o "$BINARY_NAME" .
    
    mv "$BINARY_NAME" /tmp/
    cd /
    rm -rf "$TEMP_DIR"
}

# Add ~/bin to PATH in shell config
add_to_path() {
    SHELL_RC=""
    if [ -f "$HOME/.zshrc" ]; then
        SHELL_RC="$HOME/.zshrc"
    elif [ -f "$HOME/.bashrc" ]; then
        SHELL_RC="$HOME/.bashrc"
    elif [ -f "$HOME/.profile" ]; then
        SHELL_RC="$HOME/.profile"
    fi

    if [ -n "$SHELL_RC" ]; then
        if ! grep -q 'export PATH="$HOME/bin:$PATH"' "$SHELL_RC" 2>/dev/null; then
            echo '' >> "$SHELL_RC"
            echo '# Kite CLI' >> "$SHELL_RC"
            echo 'export PATH="$HOME/bin:$PATH"' >> "$SHELL_RC"
            info "Added ~/bin to PATH in $SHELL_RC"
        fi
    fi
}

# Setup shell completions
setup_completions() {
    # Only setup zsh completions if using zsh
    if [ -f "$HOME/.zshrc" ]; then
        ZSHRC="$HOME/.zshrc"
        COMPLETIONS_DIR="$HOME/.zsh/completions"
        
        # Create completions directory
        mkdir -p "$COMPLETIONS_DIR"
        
        # Generate completion file
        "${INSTALL_DIR}/${BINARY_NAME}" completion zsh > "$COMPLETIONS_DIR/_kite" 2>/dev/null
        if [ $? -eq 0 ]; then
            info "Generated zsh completions"
        fi
        
        # Add fpath if not already present
        if ! grep -q 'fpath=(~/.zsh/completions $fpath)' "$ZSHRC" 2>/dev/null; then
            echo '' >> "$ZSHRC"
            echo '# Zsh completions' >> "$ZSHRC"
            echo 'fpath=(~/.zsh/completions $fpath)' >> "$ZSHRC"
            info "Added fpath to $ZSHRC"
        fi
        
        # Add compinit if not already present
        if ! grep -q 'autoload -Uz compinit' "$ZSHRC" 2>/dev/null && ! grep -q 'autoload -U compinit' "$ZSHRC" 2>/dev/null; then
            echo 'autoload -Uz compinit' >> "$ZSHRC"
            echo 'compinit' >> "$ZSHRC"
            info "Added compinit to $ZSHRC"
        fi
    fi
}

# Install binary to ~/bin
install() {
    info "Installing to ${INSTALL_DIR}/${BINARY_NAME}..."
    
    mkdir -p "$INSTALL_DIR"
    mv /tmp/"$BINARY_NAME" "${INSTALL_DIR}/${BINARY_NAME}"
    chmod +x "${INSTALL_DIR}/${BINARY_NAME}"
    
    add_to_path
    setup_completions
}

# Verify installation
verify() {
    if [ -x "${INSTALL_DIR}/${BINARY_NAME}" ]; then
        info "Installation successful!"
        echo ""
        echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
        echo ""
        echo "  Kite CLI installed successfully!"
        echo ""
        echo "  >>> Reload your terminal or run: source ~/.zshrc"
        echo ""
        echo "  Next steps:"
        echo "    1. kite auth         # Configure API credentials"
        echo "    2. kite login        # Authenticate with Kite"
        echo "    3. kite buy ITC 10   # Start trading!"
        echo ""
        echo "  Get help: kite --help"
        echo ""
        echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    else
        error "Installation failed. Please check errors above."
    fi
}

main() {
    echo ""
    echo "  ╔═══════════════════════════════════════╗"
    echo "  ║         Kite CLI Installer            ║"
    echo "  ║   Fast trading from your terminal     ║"
    echo "  ╚═══════════════════════════════════════╝"
    echo ""

    detect_platform
    
    if ! install_binary; then
        build_from_source
    fi
    
    install
    verify
}

main "$@"
