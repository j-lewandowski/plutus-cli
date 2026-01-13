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

echo "✅ Installed! Try:"
echo "plutus --help"
