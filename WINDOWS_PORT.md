# OpenSnitch-Windows bootstrap (WinSnitch)

## Directory tree

```text
OpenSnitch-Windows/
├── daemon/
│   ├── cmd/opensnitchd-windows/main.go
│   └── platform/windows/
│       ├── procinfo/resolver_windows.go
│       ├── service/service_windows.go
│       └── wfp/interceptor_windows.go
├── ui/
├── proto/
├── driver/
│   └── README.md
├── installer/
│   └── README.md
├── common/
│   └── README.md
├── build.ps1
├── WINDOWS_PORT.md
├── README.md
└── LICENSE
```

## Go module dependency additions

`daemon/go.mod` includes existing OpenSnitch dependencies and now adds WFP bindings:

- `github.com/tailscale/wf`

## Windows daemon entrypoint

Implemented at `daemon/cmd/opensnitchd-windows/main.go` with:

- Windows-only build tag
- Config path defaulting to `%ProgramData%`
- WFP interceptor initialization
- Console/service run modes

## One-command build/install (PowerShell, elevated)

```powershell
./build.ps1 -DeveloperMode
```

This currently builds daemon/UI and leaves driver + MSI as placeholders for the next implementation stage.
