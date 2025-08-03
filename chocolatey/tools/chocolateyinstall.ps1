$ErrorActionPreference = 'Stop'

$packageName = 'lazytrack'
$toolsDir = "$(Split-Path -parent $MyInvocation.MyCommand.Definition)"
$url = 'https://github.com/master-wayne7/lazytrack/releases/download/v0.0.0/lazytrack_Windows_x86_64.zip'
$checksum = '0000000000000000000000000000000000000000000000000000000000000000'
$checksumType = 'sha256'

$packageArgs = @{
  packageName   = $packageName
  unzipLocation = $toolsDir
  url           = $url
  checksum      = $checksum
  checksumType  = $checksumType
}

Install-ChocolateyZipPackage @packageArgs

# Create shim
$binPath = Join-Path $toolsDir "lazytrack.exe"
Install-BinFile -Name "lazytrack" -Path $binPath 