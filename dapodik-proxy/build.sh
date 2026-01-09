#!/bin/bash

# Build script for Dapodik Proxy
# ==============================

set -e

echo "Building Dapodik Proxy..."
echo ""

# Get dependencies
echo "→ Downloading dependencies..."
go mod tidy

# Build for Windows (64-bit)
echo "→ Building for Windows (amd64)..."
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o dist/dapodik-proxy-windows-amd64.exe main.go

# Build for Windows (32-bit)
echo "→ Building for Windows (386)..."
GOOS=windows GOARCH=386 go build -ldflags="-s -w" -o dist/dapodik-proxy-windows-386.exe main.go

# Build for Linux (64-bit)
echo "→ Building for Linux (amd64)..."
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o dist/dapodik-proxy-linux-amd64 main.go

# Build for macOS (64-bit)
echo "→ Building for macOS (amd64)..."
GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o dist/dapodik-proxy-darwin-amd64 main.go

# Copy config
echo "→ Copying configuration..."
cp config.yaml dist/

echo ""
echo "✓ Build complete!"
echo ""
echo "Files created in dist/:"
ls -la dist/
echo ""
echo "For Windows, copy these files to the target machine:"
echo "  - dapodik-proxy-windows-amd64.exe (64-bit)"
echo "  - dapodik-proxy-windows-386.exe (32-bit)"
echo "  - config.yaml"

