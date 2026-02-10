#!/bin/bash
set -e

REPO="j-lewandowski/plutus-cli"

OS=$(uname | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

if [[ "$ARCH" == "x86_64" ]]; then
  ARCH="amd64"
elif [[ "$ARCH" == "arm64" ]]; then
  ARCH="arm64"
fi

if [[ "$OS" == "darwin" ]]; then
  OS="mac"
fi

BIN="plutus-cli-$OS-$ARCH"

URL="https://github.com/$REPO/releases/latest/download/$BIN"

echo "Downloading $BIN..."
curl -L "$URL" -o plutus

chmod +x plutus
sudo mv plutus /usr/local/bin/plutus

# Install shell completions
SHELL_NAME=$(basename "$SHELL")

case "$SHELL_NAME" in
  zsh)
    COMP_DIR="${HOME}/.zsh/completions"
    mkdir -p "$COMP_DIR"
    plutus completion zsh > "$COMP_DIR/_plutus"
    # Ensure the completions dir is in fpath
    if ! grep -q 'fpath.*\.zsh/completions' "${HOME}/.zshrc" 2>/dev/null; then
      echo 'fpath=(~/.zsh/completions $fpath)' >> "${HOME}/.zshrc"
      echo 'autoload -Uz compinit && compinit' >> "${HOME}/.zshrc"
    fi
    echo "🐚 Zsh completions installed. Restart your shell or run: source ~/.zshrc"
    ;;
  bash)
    COMP_DIR="${HOME}/.local/share/bash-completion/completions"
    mkdir -p "$COMP_DIR"
    plutus completion bash > "$COMP_DIR/plutus"
    echo "🐚 Bash completions installed. Restart your shell to activate."
    ;;
  fish)
    COMP_DIR="${HOME}/.config/fish/completions"
    mkdir -p "$COMP_DIR"
    plutus completion fish > "$COMP_DIR/plutus.fish"
    echo "🐚 Fish completions installed. Restart your shell to activate."
    ;;
  *)
    echo "ℹ️  Shell completions not auto-installed for $SHELL_NAME."
    echo "   Run 'plutus completion --help' to set them up manually."
    ;;
esac

echo ""
echo "✅ Installed! Try:"
echo "plutus --help"
