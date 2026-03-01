# OpenSnitch-Windows

# OpenSnitch-Windows Implementation Details

This is a complete, production-ready port structure of OpenSnitch to Windows.

## Windows Specifics

1. **Service Integration:** The Daemon utilizes `golang.org/x/sys/windows/svc` to natively run as a `LocalSystem` service.
2. **Network Interception:**
   - A kernel-mode driver (`driver/OpenSnitchCallout.c`) built with the Windows Filtering Platform (WFP) acts as the connection suspended gateway.
   - It intercepts `ALE_AUTH_CONNECT_V4` and `V6` to block processes while asking the `OpenSnitch` daemon.
3. **Process Path Extraction:** It resolves processes correctly using `QueryFullProcessImageNameW` directly from Go user-mode without fragile toolhelp32 parsing.
4. **User Tracking:** Translates tokens directly to Windows SIDs.
5. **Python UI Adaptation:** We inject `windows_hacks.py` to overwrite `XDG_DATA_HOME` variables into Windows `%ProgramData%\OpenSnitch-Windows` to properly manage JSON configurations and rules.

## Build and Install (One Command)
Requires: Go 1.23, Python 3.12+, WiX Toolset.

```powershell
.\build.ps1 -BuildDaemon -BuildUI
# Use WiX Candle/Light to compile installer/OpenSnitch.wxs
```

## Known Limitations
* **Code Signing:** To run the full KMDF driver on modern Windows (10/11), it must be signed with an EV Code Signing certificate from Microsoft. You can use "Test Mode" (`bcdedit /set testsigning on`) for development.
