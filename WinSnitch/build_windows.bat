@echo off
setlocal
cd /d %~dp0

echo [1/2] Building WinSnitch binaries...
powershell -ExecutionPolicy Bypass -File .\build.ps1
if %errorlevel% neq 0 (
  echo Build failed.
  exit /b %errorlevel%
)

echo [2/2] Creating quick-start EXE/BAT with Python...
py -3 .\tools\package_with_python.py
if %errorlevel% neq 0 (
  echo Python packaging step had errors.
  exit /b %errorlevel%
)

echo Done.
echo Output folder: %cd%\dist
endlocal
