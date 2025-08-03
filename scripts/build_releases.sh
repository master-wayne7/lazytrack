#!/bin/bash

# Build script for LazyTrack cross-platform releases
echo "🚀 Building LazyTrack for all platforms..."

# Create releases directory
mkdir -p releases

# Build for different platforms
echo "📦 Building for Windows..."
GOOS=windows GOARCH=amd64 go build -o releases/lazytrack-windows-amd64.exe
GOOS=windows GOARCH=386 go build -o releases/lazytrack-windows-386.exe

echo "📦 Building for macOS..."
GOOS=darwin GOARCH=amd64 go build -o releases/lazytrack-macos-amd64
GOOS=darwin GOARCH=arm64 go build -o releases/lazytrack-macos-arm64

echo "📦 Building for Linux..."
GOOS=linux GOARCH=amd64 go build -o releases/lazytrack-linux-amd64
GOOS=linux GOARCH=386 go build -o releases/lazytrack-linux-386
GOOS=linux GOARCH=arm64 go build -o releases/lazytrack-linux-arm64

echo "📦 Building for current platform..."
go build -o releases/lazytrack

# Create checksums
echo "🔍 Creating checksums..."
cd releases
sha256sum * > checksums.txt
cd ..

echo "✅ Build complete! Files created in releases/ directory:"
ls -la releases/

echo ""
echo "📋 Next steps:"
echo "1. Create a GitHub release"
echo "2. Upload all files from releases/ directory"
echo "3. Add installation instructions to README" 