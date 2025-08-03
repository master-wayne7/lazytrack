@echo off
echo 🚀 Building LazyTrack for all platforms...

REM Create releases directory
if not exist releases mkdir releases

REM Build for different platforms
echo 📦 Building for Windows...
set GOOS=windows
set GOARCH=amd64
go build -o releases/lazytrack-windows-amd64.exe

set GOOS=windows
set GOARCH=386
go build -o releases/lazytrack-windows-386.exe

echo 📦 Building for macOS...
set GOOS=darwin
set GOARCH=amd64
go build -o releases/lazytrack-macos-amd64

set GOOS=darwin
set GOARCH=arm64
go build -o releases/lazytrack-macos-arm64

echo 📦 Building for Linux...
set GOOS=linux
set GOARCH=amd64
go build -o releases/lazytrack-linux-amd64

set GOOS=linux
set GOARCH=386
go build -o releases/lazytrack-linux-386

set GOOS=linux
set GOARCH=arm64
go build -o releases/lazytrack-linux-arm64

echo 📦 Building for current platform...
go build -o releases/lazytrack.exe

echo ✅ Build complete! Files created in releases/ directory:
dir releases

echo.
echo 📋 Next steps:
echo 1. Create a GitHub release
echo 2. Upload all files from releases/ directory
echo 3. Add installation instructions to README

pause 