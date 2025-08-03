#!/bin/bash

# LazyTrack Installation Script
echo "🚀 Installing LazyTrack..."

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

# Map architecture
case $ARCH in
    x86_64) ARCH="amd64" ;;
    aarch64) ARCH="arm64" ;;
    i386) ARCH="386" ;;
    *) echo "❌ Unsupported architecture: $ARCH"; exit 1 ;;
esac

# Map OS
case $OS in
    darwin) OS="darwin" ;;
    linux) OS="linux" ;;
    *) echo "❌ Unsupported OS: $OS"; exit 1 ;;
esac

# Download URL
BINARY_NAME="lazytrack-$OS-$ARCH"
if [ "$OS" = "darwin" ]; then
    BINARY_NAME="lazytrack-darwin-$ARCH"
fi

DOWNLOAD_URL="https://github.com/master-wayne7/lazytrack/releases/latest/download/$BINARY_NAME"

echo "📦 Downloading LazyTrack for $OS/$ARCH..."
echo "🔗 URL: $DOWNLOAD_URL"

# Download the binary
if command -v curl >/dev/null 2>&1; then
    curl -L -o lazytrack "$DOWNLOAD_URL"
elif command -v wget >/dev/null 2>&1; then
    wget -O lazytrack "$DOWNLOAD_URL"
else
    echo "❌ Neither curl nor wget found. Please install one of them."
    exit 1
fi

# Make executable
chmod +x lazytrack

# Install to system PATH
if [ "$EUID" -eq 0 ]; then
    # Running as root
    mv lazytrack /usr/local/bin/
    echo "✅ LazyTrack installed to /usr/local/bin/lazytrack"
else
    # Not running as root, ask user
    echo ""
    echo "🤔 Install to system PATH? (requires sudo)"
    read -p "Install to /usr/local/bin/ [y/N]: " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        sudo mv lazytrack /usr/local/bin/
        echo "✅ LazyTrack installed to /usr/local/bin/lazytrack"
    else
        echo "📁 LazyTrack downloaded to current directory as 'lazytrack'"
        echo "💡 You can run it with: ./lazytrack"
    fi
fi

echo ""
echo "🎉 Installation complete!"
echo ""
echo "📖 Usage examples:"
echo "  lazytrack code 2h          # Log 2 hours of coding"
echo "  lazytrack water 8x         # Log 8 glasses of water"
echo "  lazytrack summary          # View your progress"
echo "  lazytrack daemon           # Run automatic reminders"
echo ""
echo "📚 For more information, visit: https://github.com/master-wayne7/lazytrack" 