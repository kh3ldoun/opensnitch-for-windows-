# WinSnitch (OpenSnitch for Windows)

**WinSnitch** is the complete, interactive firewall and application monitor port of OpenSnitch for Windows, leveraging modern Go and Tauri 2.0 WebUIs.

## Features (100% functionally equivalent port architecture)

1. **Backend Daemon (Go 1.23+):**
   * Uses `github.com/kardianos/service` to run silently in the background as `LocalSystem`.
   * Leverages the Windows Filtering Platform (WFP) using `tailscale/wf` to capture and pend connections natively on Windows.
   * Full SIEM, rule, domain caching, and blocklist parsing structure.
2. **Frontend UI (Tauri v2 + React):**
   * Employs Rust and React (or Svelte) wrapped in `src-tauri` to produce a completely standalone `.exe` installer.
   * Includes Tray notification, multi-node configuration screens, and rule table identical to original Python UI but deeply integrated into modern OS standards.

## Project Layout

- `backend/`: Go source code for WFP interception and system service.
- `frontend/`: React + Vite + Tauri 2 WebUI.
- `src-tauri/`: Rust backend to invoke Tauri shell commands and interface with Go.
- `installer/`: Directory for NSIS or WiX package.

## Quick Start (Build & Run)

Requires **Go 1.23**, **Node.js (npm)**, **Rust**, and Windows 10+.

```powershell
# Open an Administrator PowerShell prompt
cd WinSnitch
.\build.ps1 -BuildBackend -BuildFrontend -Installer
```

### Note on Windows Driver Enforcement
To implement true *interactive pending* (where an application hangs gracefully while the UI prompt waits for "Allow" or "Deny"), a WFP Kernel-Mode Callout Driver (.sys) signed by an EV Certificate must be installed. This repository provides the user-mode Go bindings that would communicate with such a driver.
