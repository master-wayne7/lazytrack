$ErrorActionPreference = 'Stop'

$packageName = 'lazytrack'
$toolsDir = "$(Split-Path -parent $MyInvocation.MyCommand.Definition)"
$url = 'https://github.com/master-wayne7/lazytrack/releases/download/v1.0.0/lazytrack_Windows_x86_64.zip'
$checksum = '37F8152842697A85E0B244695B4199564CD9D9A2CE10AD11D964B245209A4466'
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