$ErrorActionPreference = 'Stop'

# Remove shim
Uninstall-BinFile -Name "lazytrack"

# Remove files
$toolsDir = "$(Split-Path -parent $MyInvocation.MyCommand.Definition)"
$binPath = Join-Path $toolsDir "lazytrack.exe"

if (Test-Path $binPath) {
    Remove-Item $binPath -Force
} 