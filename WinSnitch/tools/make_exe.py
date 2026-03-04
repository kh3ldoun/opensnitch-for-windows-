#!/usr/bin/env python3
"""WinSnitch one-click Windows build helper.

Usage (from CMD/PowerShell):
  py -3 tools\make_exe.py --mode portable
  py -3 tools\make_exe.py --mode installer
"""

from __future__ import annotations

import argparse
import os
import shutil
import subprocess
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
DIST = ROOT / "dist"


class BuildError(RuntimeError):
    pass


def run(cmd: list[str], cwd: Path) -> None:
    print(f"\n>> {' '.join(cmd)}")
    completed = subprocess.run(cmd, cwd=str(cwd), check=False)
    if completed.returncode != 0:
        raise BuildError(f"Command failed ({completed.returncode}): {' '.join(cmd)}")


def require(bin_name: str, hint: str) -> None:
    if shutil.which(bin_name):
        return
    raise BuildError(f"Missing required tool: {bin_name}. {hint}")


def ensure_windows() -> None:
    if os.name != "nt":
        raise BuildError("This script must run on Windows.")


def build_backend() -> None:
    require("go", "Install Go 1.23+ from https://go.dev/dl/")
    backend = ROOT / "backend"
    run(["go", "mod", "tidy"], backend)
    DIST.mkdir(parents=True, exist_ok=True)
    run(["go", "build", "-o", str(DIST / "winsnitchd.exe"), "./cmd/winsnitchd"], backend)


def build_frontend() -> None:
    require("npm", "Install Node.js 20+ from https://nodejs.org/")
    frontend = ROOT / "frontend"
    lock = frontend / "package-lock.json"
    if lock.exists():
        run(["npm", "ci"], frontend)
    else:
        run(["npm", "install"], frontend)
    run(["npm", "run", "build"], frontend)


def build_tauri() -> None:
    require("cargo", "Install Rust (MSVC toolchain) from https://rustup.rs/")
    require("cargo-tauri", "Install with: cargo install tauri-cli --version '^2'")
    tauri = ROOT / "src-tauri"
    run(["cargo", "tauri", "build"], tauri)


def collect_artifacts(mode: str) -> None:
    DIST.mkdir(parents=True, exist_ok=True)
    portable = ROOT / "src-tauri" / "target" / "release" / "winsnitch.exe"
    if portable.exists():
        shutil.copy2(portable, DIST / "WinSnitch.exe")

    if mode == "installer":
        bundle_dir = ROOT / "src-tauri" / "target" / "release" / "bundle" / "nsis"
        installers = sorted(bundle_dir.glob("*.exe"))
        if not installers:
            raise BuildError("NSIS installer not found. Check Tauri bundle output.")
        shutil.copy2(installers[-1], DIST / "WinSnitch-Installer.exe")


def main() -> int:
    parser = argparse.ArgumentParser(description="Build WinSnitch .exe artifacts")
    parser.add_argument("--mode", choices=["portable", "installer"], default="portable")
    args = parser.parse_args()

    try:
        ensure_windows()
        build_backend()
        build_frontend()
        build_tauri()
        collect_artifacts(args.mode)
    except BuildError as exc:
        print(f"\n[ERROR] {exc}")
        return 1

    print("\nBuild complete.")
    print(f"Artifacts in: {DIST}")
    print("- WinSnitch.exe")
    if args.mode == "installer":
        print("- WinSnitch-Installer.exe")
    print("- winsnitchd.exe")
    return 0


if __name__ == "__main__":
    sys.exit(main())
