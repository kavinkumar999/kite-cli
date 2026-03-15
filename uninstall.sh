#!/bin/bash
set -e

# Kite CLI Uninstaller
# Usage: curl -sSL https://raw.githubusercontent.com/kavinkumar999/kite-cli/main/uninstall.sh | bash

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

info() { echo -e "${GREEN}[INFO]${NC} $1"; }
warn() { echo -e "${YELLOW}[WARN]${NC} $1"; }

echo ""
echo "  ╔═══════════════════════════════════════╗"
echo "  ║       Kite CLI Uninstaller            ║"
echo "  ╚═══════════════════════════════════════╝"
echo ""

# Remove binary from ~/bin
if [ -f "$HOME/bin/kite" ]; then
    rm -f "$HOME/bin/kite"
    info "Removed binary: ~/bin/kite"
else
    warn "Binary not found: ~/bin/kite"
fi

# Remove config file
if [ -f "$HOME/.kite.yaml" ]; then
    read -p "Remove config file ~/.kite.yaml? This contains your credentials. (y/N): " confirm
    if [[ "$confirm" =~ ^[Yy]$ ]]; then
        rm -f "$HOME/.kite.yaml"
        info "Removed config: ~/.kite.yaml"
    else
        warn "Kept config file: ~/.kite.yaml"
    fi
else
    warn "Config not found: ~/.kite.yaml"
fi

# Remove zsh completions
if [ -f "$HOME/.zsh/completions/_kite" ]; then
    rm -f "$HOME/.zsh/completions/_kite"
    info "Removed zsh completion: ~/.zsh/completions/_kite"
    
    # Remove empty completions directory
    if [ -d "$HOME/.zsh/completions" ] && [ -z "$(ls -A "$HOME/.zsh/completions")" ]; then
        rmdir "$HOME/.zsh/completions"
        info "Removed empty directory: ~/.zsh/completions"
    fi
    
    # Remove empty .zsh directory
    if [ -d "$HOME/.zsh" ] && [ -z "$(ls -A "$HOME/.zsh")" ]; then
        rmdir "$HOME/.zsh"
        info "Removed empty directory: ~/.zsh"
    fi
fi

# Remove bash completions
if [ -f "/etc/bash_completion.d/kite" ]; then
    sudo rm -f "/etc/bash_completion.d/kite" 2>/dev/null || true
    info "Removed bash completion"
fi

# Clean up any leftover temp files from install
rm -f /tmp/kite /tmp/kite.tar.gz 2>/dev/null || true
rm -rf /tmp/kite-cli-* 2>/dev/null || true
info "Cleaned up temp files"

# Check if ~/bin is empty and offer to remove
if [ -d "$HOME/bin" ] && [ -z "$(ls -A "$HOME/bin")" ]; then
    read -p "~/bin directory is empty. Remove it? (y/N): " confirm
    if [[ "$confirm" =~ ^[Yy]$ ]]; then
        rmdir "$HOME/bin"
        info "Removed empty directory: ~/bin"
    fi
fi

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "  Kite CLI uninstalled successfully!"
echo ""
echo "  Optional cleanup (manual):"
echo "    - Remove 'export PATH=\"\$HOME/bin:\$PATH\"' from ~/.zshrc"
echo "    - Remove 'fpath=(~/.zsh/completions \$fpath)' from ~/.zshrc"
echo "    - Remove '# Kite CLI' comment from ~/.zshrc"
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
