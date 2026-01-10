#!/bin/bash

# Build script for Dapodik Proxy
# ==============================

set -e

echo "╔══════════════════════════════════════════════════════════════╗"
echo "║           BUILDING DAPODIK PROXY                             ║"
echo "╚══════════════════════════════════════════════════════════════╝"
echo ""

# Create dist directory if not exists
mkdir -p dist

# Get dependencies for main proxy
echo "→ Downloading proxy dependencies..."
go mod tidy

# Get dependencies for installer
echo "→ Downloading installer dependencies..."
cd installer
go mod tidy
cd ..

echo ""
echo "Building Proxy Executables..."
echo "───────────────────────────────────────────────────────────────"

# Build Proxy for Windows (64-bit)
echo "→ Proxy: Windows (amd64)..."
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o dist/dapodik-proxy-windows-amd64.exe main.go

# Build Proxy for Windows (32-bit)
echo "→ Proxy: Windows (386)..."
GOOS=windows GOARCH=386 go build -ldflags="-s -w" -o dist/dapodik-proxy-windows-386.exe main.go

# Build Proxy for Linux (64-bit)
echo "→ Proxy: Linux (amd64)..."
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o dist/dapodik-proxy-linux-amd64 main.go

# Build Proxy for macOS (64-bit Intel)
echo "→ Proxy: macOS (amd64)..."
GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o dist/dapodik-proxy-darwin-amd64 main.go

# Build Proxy for macOS (ARM64 Apple Silicon)
echo "→ Proxy: macOS (arm64)..."
GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o dist/dapodik-proxy-darwin-arm64 main.go

echo ""
echo "Building Installer Executables..."
echo "───────────────────────────────────────────────────────────────"

# Build Installer for Windows (64-bit)
echo "→ Installer: Windows (amd64)..."
cd installer
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o ../dist/installer-windows-amd64.exe main.go

# Build Installer for Windows (32-bit)
echo "→ Installer: Windows (386)..."
GOOS=windows GOARCH=386 go build -ldflags="-s -w" -o ../dist/installer-windows-386.exe main.go
cd ..

echo ""
echo "Copying configuration files..."
echo "───────────────────────────────────────────────────────────────"
cp config.yaml dist/
cp README.md dist/

echo ""
echo "Creating distribution packages..."
echo "───────────────────────────────────────────────────────────────"

cd dist

# Windows 64-bit package
echo "→ Package: Windows (amd64)..."
rm -rf dapodik-proxy-windows-amd64
mkdir -p dapodik-proxy-windows-amd64
cp dapodik-proxy-windows-amd64.exe dapodik-proxy-windows-amd64/
cp installer-windows-amd64.exe dapodik-proxy-windows-amd64/installer.exe
cp config.yaml dapodik-proxy-windows-amd64/
cp README.md dapodik-proxy-windows-amd64/
zip -q -r dapodik-proxy-windows-amd64.zip dapodik-proxy-windows-amd64/

# Windows 32-bit package
echo "→ Package: Windows (386)..."
rm -rf dapodik-proxy-windows-386
mkdir -p dapodik-proxy-windows-386
cp dapodik-proxy-windows-386.exe dapodik-proxy-windows-386/
cp installer-windows-386.exe dapodik-proxy-windows-386/installer.exe
cp config.yaml dapodik-proxy-windows-386/
cp README.md dapodik-proxy-windows-386/
zip -q -r dapodik-proxy-windows-386.zip dapodik-proxy-windows-386/

# Linux package
echo "→ Package: Linux (amd64)..."
rm -rf dapodik-proxy-linux-amd64-pkg
mkdir -p dapodik-proxy-linux-amd64-pkg
cp dapodik-proxy-linux-amd64 dapodik-proxy-linux-amd64-pkg/dapodik-proxy
cp config.yaml dapodik-proxy-linux-amd64-pkg/
cp README.md dapodik-proxy-linux-amd64-pkg/
tar -czf dapodik-proxy-linux-amd64.tar.gz dapodik-proxy-linux-amd64-pkg/

cd ..

echo ""
echo "╔══════════════════════════════════════════════════════════════╗"
echo "║           ✓ BUILD COMPLETE!                                  ║"
echo "╚══════════════════════════════════════════════════════════════╝"
echo ""
echo "Distribution Packages:"
echo "  • dapodik-proxy-windows-amd64.zip (64-bit Windows + Installer)"
echo "  • dapodik-proxy-windows-386.zip (32-bit Windows + Installer)"
echo "  • dapodik-proxy-linux-amd64.tar.gz (Linux 64-bit)"
echo ""
echo "Installation Instructions:"
echo "  Windows:"
echo "    1. Extract the ZIP file"
echo "    2. Right-click installer.exe → 'Run as Administrator'"
echo "    3. Select option [1] to install service"
echo ""
echo "  Linux:"
echo "    1. Extract: tar -xzf dapodik-proxy-linux-amd64.tar.gz"
echo "    2. Run: ./dapodik-proxy-linux-amd64"
echo ""

