#!/bin/sh
# install.sh — install bb (Bitbucket Cloud CLI)
# Usage: curl -sSL https://tyrantkhan.github.io/bitbucket-cli/install.sh | sh
set -eu

REPO="tyrantkhan/bitbucket-cli"
INSTALL_DIR="${HOME}/.local/bin"

# Detect OS
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
case "$OS" in
    linux)  OS="linux" ;;
    darwin) OS="darwin" ;;
    *)      echo "Unsupported OS: $OS" >&2; exit 1 ;;
esac

# Detect architecture
ARCH="$(uname -m)"
case "$ARCH" in
    x86_64|amd64)  ARCH="amd64" ;;
    arm64|aarch64) ARCH="arm64" ;;
    *)             echo "Unsupported architecture: $ARCH" >&2; exit 1 ;;
esac

# Get latest version
VERSION="$(curl -sSL "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | sed 's/.*"v\(.*\)".*/\1/')"
if [ -z "$VERSION" ]; then
    echo "Failed to determine latest version" >&2
    exit 1
fi

ARCHIVE="bb_${VERSION}_${OS}_${ARCH}.tar.gz"
URL="https://github.com/${REPO}/releases/download/v${VERSION}/${ARCHIVE}"

echo "Installing bb v${VERSION} (${OS}/${ARCH})..."

TMPDIR="$(mktemp -d)"
trap 'rm -rf "$TMPDIR"' EXIT

curl -sSL "$URL" -o "${TMPDIR}/${ARCHIVE}"
tar -xzf "${TMPDIR}/${ARCHIVE}" -C "$TMPDIR"

mkdir -p "$INSTALL_DIR"
mv "${TMPDIR}/bb" "${INSTALL_DIR}/bb"
chmod +x "${INSTALL_DIR}/bb"

echo "Installed bb to ${INSTALL_DIR}/bb"

# Check if INSTALL_DIR is in PATH
case ":${PATH}:" in
    *":${INSTALL_DIR}:"*)
        echo "Run 'bb --help' to get started."
        ;;
    *)
        echo ""
        echo "WARNING: ${INSTALL_DIR} is not in your PATH."
        echo "The 'bb' command won't be found until you add it."
        echo ""
        SHELL_NAME="$(basename "${SHELL:-sh}")"
        case "$SHELL_NAME" in
            zsh)  RC_FILE="~/.zshrc" ;;
            bash) RC_FILE="~/.bashrc" ;;
            fish) RC_FILE="~/.config/fish/config.fish" ;;
            *)    RC_FILE="your shell's config file" ;;
        esac
        if [ "$SHELL_NAME" = "fish" ]; then
            echo "Run this to add it:"
            echo "  fish_add_path ${INSTALL_DIR}"
        elif [ "$RC_FILE" = "your shell's config file" ]; then
            echo "Add this to your shell's config file (e.g. ~/.profile):"
            echo "  export PATH=\"${INSTALL_DIR}:\$PATH\""
        else
            echo "Run this to add it:"
            echo "  echo 'export PATH=\"${INSTALL_DIR}:\$PATH\"' >> ${RC_FILE}"
        fi
        echo ""
        echo "Then restart your shell or run: exec \$SHELL"
        ;;
esac
