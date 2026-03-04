# WinSnitch (OpenSnitch-Windows)

WinSnitch is a Windows-first OpenSnitch port with:
- Go Windows service daemon (`winsnitchd`) using `kardianos/service`
- WFP interception layer using `tailscale/wf` abstraction
- Tauri 2 + React + Tailwind desktop UI with tray and notifications
- JSON logs + JSON rules under `%ProgramData%\\WinSnitch`

## Directory tree

```text
WinSnitch/
├── backend/
│   ├── cmd/winsnitchd/main.go
│   ├── go.mod
│   └── internal/
│       ├── api/ws.go
│       ├── blocklist/blocklist.go
│       ├── config/config.go
│       ├── dnscache/cache.go
│       ├── events/events.go
│       ├── logging/jsonlog.go
│       ├── multinode/manager.go
│       ├── rules/rules.go
│       ├── service/service.go
│       └── wfp/{engine_windows.go,engine_stub.go}
├── frontend/
│   ├── package.json
│   ├── src/
│   │   ├── components/DecisionModal.tsx
│   │   ├── hooks/useLiveEvents.ts
│   │   ├── lib/types.ts
│   │   ├── pages/Dashboard.tsx
│   │   ├── main.tsx
│   │   └── styles/globals.css
│   └── {index.html,postcss.config.js,tailwind.config.ts,tsconfig.json,vite.config.ts}
├── installer/winsnitch.nsi
├── src-tauri/
│   ├── Cargo.toml
│   ├── tauri.conf.json
│   └── src/{lib.rs,main.rs}
└── build.ps1
```

## 2-minute build

1. Open an elevated PowerShell on Windows 10/11.
2. Install Go 1.23+, Rust stable (MSVC), Node 20+, and Tauri prerequisites.
3. Run:

```powershell
cd .\WinSnitch
.\build.ps1
```

Outputs:
- `dist\WinSnitch.exe` (portable Tauri app)
- `dist\winsnitchd.exe` (service daemon)

## Run

```powershell
cd .\WinSnitch\dist
.\winsnitchd.exe --install
sc start WinSnitch
.\WinSnitch.exe
```


## One-click EXE/BAT (for your request)

If you want a ready launcher (`.exe`) or a `.bat` generated with Python support:

```bat
cd WinSnitch
build_windows.bat
```

This does:
- build backend + Tauri app
- run `tools\package_with_python.py`
- if PyInstaller exists: creates `dist\WinSnitch-QuickStart.exe`
- otherwise: creates `dist\Run-WinSnitch.bat`

## Windows ops notes

- Defender exclusion (optional):
  `Add-MpPreference -ExclusionPath "C:\ProgramData\WinSnitch"`
- VPN: interception is applied at WFP layer and can include tunnel interfaces.
- Sleep/resume: service restarts filters on startup; add a Task Scheduler wake trigger for strict environments.
