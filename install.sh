#!/usr/bin/env sh
set -e

REPO="SOORAJTS2001/uplog"
BIN="uplog"
INSTALL_DIR="$HOME/.local/bin"


# Detect OS

OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
case "$OS" in
  linux)   OS="linux" ;;
  darwin)  OS="darwin" ;;
  *)
    echo "Unsupported OS: $OS"
    exit 1
    ;;
esac


# Detect ARCH

ARCH="$(uname -m)"
case "$ARCH" in
  x86_64|amd64) ARCH="amd64" ;;
  arm64|aarch64) ARCH="arm64" ;;
  *)
    echo "Unsupported architecture: $ARCH"
    exit 1
    ;;
esac


# Version

VERSION="${VERSION:-latest}"

if [ "$VERSION" = "latest" ]; then
  URL="https://github.com/$REPO/releases/latest/download/$BIN-$OS-$ARCH"
else
  URL="https://github.com/$REPO/releases/download/$VERSION/$BIN-$OS-$ARCH"
fi


# Install dir

mkdir -p "$INSTALL_DIR"

TMP="$(mktemp)"
echo "Downloading $BIN ($OS/$ARCH) version - $VERSION"
curl -fsSL "$URL" -o "$TMP"

chmod +x "$TMP"
mv "$TMP" "$INSTALL_DIR/$BIN"

echo "Installed $BIN to $INSTALL_DIR/$BIN"


# PATH check

if ! echo "$PATH" | tr ':' '\n' | grep -qx "$INSTALL_DIR"; then
  echo ""
  echo "$INSTALL_DIR is not in your PATH"
  echo "Add this to your shell config:"
  echo ""
  echo "  export PATH=\"\$PATH:$INSTALL_DIR\""
  echo ""
fi

echo "Run: $BIN --help"
