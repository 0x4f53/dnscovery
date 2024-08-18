#!/bin/bash

# Define the program name and version
PROGRAM_NAME="dnscovery"
VERSION="1.0.0"

# Define the operating systems and architectures you want to build for
OS_ARCH_LIST=(
    "linux/amd64"
    "linux/arm64"
    "linux/386"
    "darwin/amd64"
    "darwin/arm64"
    "windows/amd64"
    "windows/386"
)

# Create an output directory for the binaries
OUTPUT_DIR=".build/binaries/"

rm -rf $OUTPUT_DIR
rm $PROGRAM_NAME

mkdir -p "$OUTPUT_DIR"

# Iterate over each OS/Arch combination and build the binary
for OS_ARCH in "${OS_ARCH_LIST[@]}"; do
    IFS="/" read -r OS ARCH <<< "$OS_ARCH"
    
    # Set the output binary name and path
    BINARY_NAME="${PROGRAM_NAME}-${VERSION}-${OS}-${ARCH}"
    
    # Add .exe extension for Windows binaries
    if [ "$OS" = "windows" ]; then
        BINARY_NAME="${BINARY_NAME}.exe"
    fi

    OUTPUT_PATH="${OUTPUT_DIR}/${BINARY_NAME}"
    
    # Set the environment variables and build the binary
    echo "Building for $OS/$ARCH..."
    env GOOS="$OS" GOARCH="$ARCH" go build -o "$OUTPUT_PATH"
    
    if [ $? -ne 0 ]; then
        echo "Failed to build for $OS/$ARCH"
        exit 1
    fi
    
done

echo "All binaries built successfully!"