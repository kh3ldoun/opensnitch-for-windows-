param(
  [switch]$Release = $true
)

$ErrorActionPreference = 'Stop'
$root = Split-Path -Parent $MyInvocation.MyCommand.Path

Write-Host '[1/4] Building backend service...'
Push-Location "$root/backend"
go mod tidy
go build -o "$root/dist/winsnitchd.exe" ./cmd/winsnitchd
Pop-Location

Write-Host '[2/4] Building frontend...'
Push-Location "$root/frontend"
npm install
npm run build
Pop-Location

Write-Host '[3/4] Building Tauri desktop app...'
Push-Location "$root/src-tauri"
cargo tauri build
Pop-Location

Write-Host '[4/4] Collecting artifacts...'
New-Item -ItemType Directory -Force "$root/dist" | Out-Null
Copy-Item "$root/src-tauri/target/release/winsnitch.exe" "$root/dist/WinSnitch.exe" -Force
Write-Host "Done. Portable executable: $root/dist/WinSnitch.exe"
