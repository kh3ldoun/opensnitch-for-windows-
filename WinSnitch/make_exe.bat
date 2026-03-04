@echo off
setlocal

REM One-click EXE build wrapper (portable by default)
set MODE=%1
if "%MODE%"=="" set MODE=portable

where py >nul 2>&1
if errorlevel 1 (
  echo [ERROR] Python launcher 'py' was not found. Install Python 3.x first.
  exit /b 1
)

py -3 "%~dp0tools\make_exe.py" --mode %MODE%
if errorlevel 1 exit /b 1

echo.
echo Done. Output files are in "%~dp0dist"
exit /b 0
