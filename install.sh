#!/usr/bin/env sh
set -e

REPO="SOORAJTS2001/uplog"        # change if needed
BIN="uplog"
INSTALL_DIR="/usr/local/bin"

OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
case "$OS" in
  linux)   OS="linux" ;;
  darwin)  OS="darwin" ;;
  *)
    echo "Unsupported OS: $OS"
    exit 1
    ;;
esac


ARCH="$(uname -m)"
case "$ARCH" in
  x86_64|amd64) ARCH="amd64" ;;
  arm64|aarch64) ARCH="arm64" ;;
  *)
    echo "Unsupported architecture: $ARCH"
    exit 1
    ;;
esac


VERSION="${VERSION:-latest}"

if [ "$VERSION" = "latest" ]; then
  DOWNLOAD_URL="https://github.com/$REPO/releases/latest/download/$BIN-$OS-$ARCH"
else
  DOWNLOAD_URL="https://github.com/$REPO/releases/download/$VERSION/$BIN-$OS-$ARCH"
fi


TMP="$(mktemp)"
echo "Downloading $BIN for $OS/$ARCH..."
curl -fsSL "$DOWNLOAD_URL" -o "$TMP"


chmod +x "$TMP"

if [ -w "$INSTALL_DIR" ]; then
  mv "$TMP" "$INSTALL_DIR/$BIN"
else
  echo "Installing to $INSTALL_DIR (sudo required)"
  sudo mv "$TMP" "$INSTALL_DIR/$BIN"
fi

echo "$BIN installed at $INSTALL_DIR/$BIN"
echo "Run: $BIN --help"
