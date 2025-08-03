# LazyTrack Installation Script for Windows
Write-Host "🚀 Installing LazyTrack..." -ForegroundColor Green

# Detect architecture
$Architecture = if ([Environment]::Is64BitOperatingSystem) { "amd64" } else { "386" }
$BinaryName = "lazytrack-windows-$Architecture.exe"
$DownloadUrl = "https://github.com/master-wayne7/lazytrack/releases/latest/download/$BinaryName"

Write-Host "📦 Downloading LazyTrack for Windows/$Architecture..." -ForegroundColor Yellow
Write-Host "🔗 URL: $DownloadUrl" -ForegroundColor Cyan

try {
    # Download the binary
    Invoke-WebRequest -Uri $DownloadUrl -OutFile "lazytrack.exe" -UseBasicParsing
    Write-Host "✅ Download completed!" -ForegroundColor Green
}
catch {
    Write-Host "❌ Download failed: $($_.Exception.Message)" -ForegroundColor Red
    exit 1
}

# Ask user if they want to install to system PATH
Write-Host ""
Write-Host "🤔 Install to system PATH?" -ForegroundColor Yellow
$Response = Read-Host "Install to C:\Windows\System32\ [y/N]"

if ($Response -eq "y" -or $Response -eq "Y") {
    try {
        # Check if running as administrator
        if (-NOT ([Security.Principal.WindowsPrincipal] [Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole] "Administrator")) {
            Write-Host "⚠️  This operation requires administrator privileges." -ForegroundColor Yellow
            Write-Host "💡 Please run PowerShell as Administrator and try again." -ForegroundColor Cyan
            Write-Host "📁 LazyTrack downloaded to current directory as 'lazytrack.exe'" -ForegroundColor Green
            Write-Host "💡 You can run it with: .\lazytrack.exe" -ForegroundColor Cyan
        }
        else {
            # Copy to system directory
            Copy-Item "lazytrack.exe" "C:\Windows\System32\" -Force
            Write-Host "✅ LazyTrack installed to C:\Windows\System32\lazytrack.exe" -ForegroundColor Green
        }
    }
    catch {
        Write-Host "❌ Installation failed: $($_.Exception.Message)" -ForegroundColor Red
        Write-Host "📁 LazyTrack downloaded to current directory as 'lazytrack.exe'" -ForegroundColor Green
        Write-Host "💡 You can run it with: .\lazytrack.exe" -ForegroundColor Cyan
    }
}
else {
    Write-Host "📁 LazyTrack downloaded to current directory as 'lazytrack.exe'" -ForegroundColor Green
    Write-Host "💡 You can run it with: .\lazytrack.exe" -ForegroundColor Cyan
}

Write-Host ""
Write-Host "🎉 Installation complete!" -ForegroundColor Green
Write-Host ""
Write-Host "📖 Usage examples:" -ForegroundColor Yellow
Write-Host "  lazytrack code 2h          # Log 2 hours of coding" -ForegroundColor White
Write-Host "  lazytrack water 8x         # Log 8 glasses of water" -ForegroundColor White
Write-Host "  lazytrack summary          # View your progress" -ForegroundColor White
Write-Host "  lazytrack daemon           # Run automatic reminders" -ForegroundColor White
Write-Host ""
Write-Host "📚 For more information, visit: https://github.com/master-wayne7/lazytrack" -ForegroundColor Cyan 