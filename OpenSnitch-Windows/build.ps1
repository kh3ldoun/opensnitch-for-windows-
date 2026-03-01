# OpenSnitch-Windows Build Script

param (
    [switch]$BuildDaemon = $true,
    [switch]$BuildUI = $true,
    [switch]$BuildDriver = $false
)

$ErrorActionPreference = "Stop"
$WorkingDir = (Get-Item -Path ".\" -Verbose).FullName

Write-Host "--- OpenSnitch-Windows Build Process ---"

if ($BuildDaemon) {
    Write-Host "Building Go Daemon..."
    cd $WorkingDir\daemon
    go mod tidy
    # Build for Windows x64
    $env:GOOS = "windows"
    $env:GOARCH = "amd64"
    go build -o ..\bin\opensnitchd.exe main.go process.go wfp.go decision.go grpc.go
    if ($?) { Write-Host "Daemon build successful." -ForegroundColor Green }
}

if ($BuildUI) {
    Write-Host "Packaging Python UI with PyInstaller..."
    cd $WorkingDir\ui
    # In a real environment, you need Python and PyInstaller installed
    # pip install -r requirements.txt
    # pip install pyinstaller
    # pyinstaller --noconfirm --onedir --windowed --add-data "opensnitch\res;opensnitch\res" --name "OpenSnitchUI" opensnitch_ui_windows.py
    Write-Host "UI Build step complete (simulated)." -ForegroundColor Green
}

if ($BuildDriver) {
    Write-Host "Building WFP Kernel Driver..."
    cd $WorkingDir\driver
    # Requires Visual Studio Build Tools, WDK, and MSBuild
    # msbuild OpenSnitchCallout.sln /p:Configuration=Release /p:Platform=x64
    Write-Host "Driver Build step complete (simulated)." -ForegroundColor Green
}

Write-Host "--- Build Finished ---"
