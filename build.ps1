param(
  [switch]$DeveloperMode
)

$ErrorActionPreference = 'Stop'

Write-Host '==> Building Go daemon'
Push-Location daemon
$env:GOOS='windows'
$env:GOARCH='amd64'
go build -o ..\bin\opensnitchd.exe .\cmd\opensnitchd-windows
Pop-Location

Write-Host '==> Building UI package'
python -m pip install -r ui/requirements.txt
python -m PyInstaller ui/opensnitch-ui.spec

Write-Host '==> Driver build (placeholder)'
if ($DeveloperMode) {
  Write-Host 'Developer mode enabled: test-signing workflow expected.'
}

Write-Host 'Build complete.'
