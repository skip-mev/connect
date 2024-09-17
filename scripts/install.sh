#!/bin/bash

set -e

CONNECT_RELEASES_URL="https://api.github.com/repos/skip-mev/connect/releases/latest"

# Determine the system architecture
ARCH=$(uname -m)
OS=$(uname -s | tr '[:upper:]' '[:lower:]')

# Fetch the latest release information
echo "Fetching latest release information..."
RELEASE_INFO=$(curl -Ls ${CONNECT_RELEASES_URL})
VERSION=$(echo "${RELEASE_INFO}" | grep -o '"tag_name": "v[^"]*' | cut -d'"' -f4)
VERSION=${VERSION#v}  # Remove the 'v' prefix

# Map architecture to release file name
case "${ARCH}" in
    x86_64)
        if [ "${OS}" = "darwin" ]; then
            FILE_NAME="connect-${VERSION}-darwin-amd64.tar.gz"
        else
            FILE_NAME="connect-${VERSION}-linux-amd64.tar.gz"
        fi
        ;;
    aarch64|arm64)
        if [ "${OS}" = "darwin" ]; then
            FILE_NAME="connect-${VERSION}-darwin-arm64.tar.gz"
        else
            FILE_NAME="connect-${VERSION}-linux-arm64.tar.gz"
        fi
        ;;
    i386|i686)
        FILE_NAME="connect-${VERSION}-linux-386.tar.gz"
        ;;
    *)
        echo "Unsupported architecture: ${ARCH}"
        exit 1
        ;;
esac

# Get download URL for the specific file
DOWNLOAD_URL=$(echo "${RELEASE_INFO}" | grep -o "\"browser_download_url\": \"[^\"]*${FILE_NAME}\"" | cut -d'"' -f4)

if [ -z "${DOWNLOAD_URL}" ]; then
    echo "Failed to find download URL for ${FILE_NAME}"
    exit 1
fi

# Download the release
echo "Downloading ${FILE_NAME}..."
curl -LO "${DOWNLOAD_URL}"

# Create a temporary directory for extraction
TEMP_DIR=$(mktemp -d)
echo "Extracting connect binary to ${TEMP_DIR}..."
tar -xzf "${FILE_NAME}" -C "${TEMP_DIR}"

# Find the connect binary
CONNECT_BIN=$(find "${TEMP_DIR}" -type f -name "connect")

if [ -z "${CONNECT_BIN}" ]; then
    echo "Failed to find connect binary in the extracted files"
    rm -rf "${TEMP_DIR}"
    rm "${FILE_NAME}"
    exit 1
fi

# Move the binary to /usr/local/bin
echo "Installing connect to /usr/local/bin..."
sudo mv "${CONNECT_BIN}" /usr/local/bin/

# Make it executable
sudo chmod +x /usr/local/bin/connect

# Clean up
rm -rf "${TEMP_DIR}"
rm "${FILE_NAME}"

echo "Connect ${VERSION} has been installed successfully!"
