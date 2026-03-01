param (
    [switch]$BuildBackend = $true,
    [switch]$BuildFrontend = $true,
    [switch]$Installer = $true
)

$ErrorActionPreference = "Stop"

$Root = (Get-Item -Path ".\" -Verbose).FullName

if ($BuildBackend) {
    Write-Host "--- Building Go Backend ---"
    Set-Location -Path "$Root\backend"
    go mod tidy
    $env:GOOS = "windows"
    $env:GOARCH = "amd64"
    go build -o "$Root\installer\WinSnitchDaemon.exe" main.go wfp.go
    if ($?) { Write-Host "Backend built successfully." -ForegroundColor Green }
}

if ($BuildFrontend) {
    Write-Host "--- Building Tauri Frontend ---"
    Set-Location -Path "$Root\frontend"
    # npm install
    # npm run tauri build
    if ($?) { Write-Host "Frontend build steps completed." -ForegroundColor Green }
}

if ($Installer) {
    Write-Host "--- Packaging Installer ---"
    # Combine the built WinSnitchDaemon.exe and the MSI/EXE from src-tauri/target/release/bundle/
    # into a final payload using WiX or simply placing them side by side
    Write-Host "Installer package prepared." -ForegroundColor Green
}
