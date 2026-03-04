#!/usr/bin/env python3
"""Build a one-file Windows launcher EXE (via PyInstaller) for WinSnitch.

Usage (from WinSnitch root on Windows):
  py tools\package_with_python.py
"""
from __future__ import annotations

import shutil
import subprocess
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
DIST = ROOT / "dist"
BUILD_DIR = ROOT / ".pybuild"

LAUNCHER_SOURCE = r'''
from __future__ import annotations
import subprocess
from pathlib import Path

ROOT = Path(__file__).resolve().parent
SERVICE_BIN = ROOT / "winsnitchd.exe"
UI_BIN = ROOT / "WinSnitch.exe"


def run(cmd: list[str]) -> int:
    proc = subprocess.run(cmd, capture_output=True, text=True)
    if proc.returncode != 0 and proc.stderr:
        print(proc.stderr.strip())
    return proc.returncode


def main() -> int:
    if not SERVICE_BIN.exists() or not UI_BIN.exists():
        print("Missing required files: winsnitchd.exe / WinSnitch.exe")
        print("Run build.ps1 first.")
        return 1

    # Try to install/start service if not installed yet.
    _ = run([str(SERVICE_BIN), "--install"])
    _ = run(["sc", "start", "WinSnitch"])

    return subprocess.call([str(UI_BIN)])


if __name__ == "__main__":
    raise SystemExit(main())
'''


def write_fallback_bat() -> Path:
    target = DIST / "Run-WinSnitch.bat"
    target.write_text(
        "@echo off\r\n"
        "setlocal\r\n"
        "cd /d %~dp0\r\n"
        "winsnitchd.exe --install\r\n"
        "sc start WinSnitch\r\n"
        "start \"\" WinSnitch.exe\r\n"
        "endlocal\r\n",
        encoding="utf-8",
    )
    return target


def main() -> int:
    DIST.mkdir(parents=True, exist_ok=True)
    launcher = BUILD_DIR / "launcher.py"
    BUILD_DIR.mkdir(parents=True, exist_ok=True)
    launcher.write_text(LAUNCHER_SOURCE, encoding="utf-8")

    pyinstaller = shutil.which("pyinstaller")
    if not pyinstaller:
        bat = write_fallback_bat()
        print(f"PyInstaller not found. Created BAT fallback: {bat}")
        return 0

    cmd = [
        pyinstaller,
        "--noconfirm",
        "--clean",
        "--onefile",
        "--noconsole",
        "--name",
        "WinSnitch-QuickStart",
        "--distpath",
        str(DIST),
        "--workpath",
        str(BUILD_DIR / "work"),
        "--specpath",
        str(BUILD_DIR),
        str(launcher),
    ]
    print("Running:", " ".join(cmd))
    res = subprocess.run(cmd)
    if res.returncode != 0:
        bat = write_fallback_bat()
        print(f"PyInstaller failed. Created BAT fallback: {bat}")
        return res.returncode

    print(f"Created launcher EXE: {DIST / 'WinSnitch-QuickStart.exe'}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
