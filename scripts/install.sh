#!/usr/bin/env bash
set -euo pipefail

REPO="${LSKV_REPO:-zemelkajakub/lskv}"
BIN_NAME="lskv"
INSTALL_DIR="/usr/local/bin"

usage() {
  cat <<'EOF'
Install lskv from the latest GitHub release.

Usage:
  install.sh [--repo owner/name] [--install-dir /path]

Options:
  --repo         GitHub repository (default: zemelkajakub/lskv)
  --install-dir  Target directory (default: /usr/local/bin)
  -h, --help     Show this help

Examples:
  curl -fsSL https://raw.githubusercontent.com/zemelkajakub/lskv/main/scripts/install.sh | sudo bash
  curl -fsSL https://raw.githubusercontent.com/zemelkajakub/lskv/main/scripts/install.sh | bash -s -- --install-dir "$HOME/.local/bin"
EOF
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --repo)
      REPO="$2"
      shift 2
      ;;
    --install-dir)
      INSTALL_DIR="$2"
      shift 2
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    *)
      echo "Unknown argument: $1" >&2
      usage
      exit 1
      ;;
  esac
done

OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

case "$OS" in
  linux|darwin)
    ;;
  *)
    echo "Unsupported OS: $OS" >&2
    exit 1
    ;;
esac

case "$ARCH" in
  x86_64|amd64)
    ARCH="amd64"
    ;;
  arm64|aarch64)
    ARCH="arm64"
    ;;
  *)
    echo "Unsupported architecture: $ARCH" >&2
    exit 1
    ;;
esac

ASSET="${BIN_NAME}_${OS}_${ARCH}.tar.gz"
URL="https://github.com/${REPO}/releases/latest/download/${ASSET}"

TMP_DIR="$(mktemp -d)"
trap 'rm -rf "$TMP_DIR"' EXIT

echo "Downloading ${ASSET} from ${REPO}..."
curl -fL "$URL" -o "$TMP_DIR/$ASSET"

tar -xzf "$TMP_DIR/$ASSET" -C "$TMP_DIR"

if [[ ! -f "$TMP_DIR/$BIN_NAME" ]]; then
  echo "Binary not found in release archive: $ASSET" >&2
  exit 1
fi

chmod +x "$TMP_DIR/$BIN_NAME"

if [[ -d "$INSTALL_DIR" ]]; then
  :
elif [[ -w "$(dirname "$INSTALL_DIR")" ]]; then
  mkdir -p "$INSTALL_DIR"
else
  echo "Creating ${INSTALL_DIR} requires elevated permissions..."
  sudo mkdir -p "$INSTALL_DIR"
fi

if [[ -w "$INSTALL_DIR" ]]; then
  install -m 0755 "$TMP_DIR/$BIN_NAME" "$INSTALL_DIR/$BIN_NAME"
else
  echo "Installing to ${INSTALL_DIR} requires elevated permissions..."
  sudo install -m 0755 "$TMP_DIR/$BIN_NAME" "$INSTALL_DIR/$BIN_NAME"
fi

echo "Installed: ${INSTALL_DIR}/${BIN_NAME}"
echo "Run: ${BIN_NAME} --help"
