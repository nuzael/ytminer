#!/bin/bash

# YTMiner Installation Script

set -e

echo "Ì∫Ä YTMiner Installation Script"
echo "==============================="

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case $ARCH in
    x86_64) ARCH="amd64" ;;
    arm64|aarch64) ARCH="arm64" ;;
    *) echo "‚ùå Unsupported architecture: $ARCH"; exit 1 ;;
esac

echo "Ì≥ã Detected: $OS $ARCH"

# Set download URL
LATEST_RELEASE=$(curl -s https://api.github.com/repos/nuzael/ytminer/releases/latest | grep "tag_name" | cut -d '"' -f 4)
BINARY_NAME="ytminer-$OS-$ARCH"

if [[ "$OS" == "mingw"* ]] || [[ "$OS" == "msys"* ]] || [[ "$OS" == "cygwin"* ]]; then
    BINARY_NAME="ytminer-windows-$ARCH.exe"
    OS="windows"
fi

echo "Ì≥¶ Downloading YTMiner $LATEST_RELEASE..."

# Create temp directory
TEMP_DIR=$(mktemp -d)
cd "$TEMP_DIR"

# Download binary
curl -L "https://github.com/nuzael/ytminer/releases/download/$LATEST_RELEASE/$BINARY_NAME" -o ytminer

# Make executable
chmod +x ytminer

# Install to /usr/local/bin (requires sudo)
if command -v sudo >/dev/null 2>&1; then
    echo "Ì¥ß Installing to /usr/local/bin (requires sudo)..."
    sudo mv ytminer /usr/local/bin/
    echo "‚úÖ YTMiner installed successfully!"
    echo "ÌæØ Run 'ytminer --help' to get started"
else
    echo "‚ö†Ô∏è  Sudo not available. Binary downloaded to: $TEMP_DIR/ytminer"
    echo "Ì¥ß Please move it to a directory in your PATH manually"
fi

# Cleanup
cd /
rm -rf "$TEMP_DIR"
