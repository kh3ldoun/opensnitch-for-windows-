OutFile "WinSnitch-Installer.exe"
InstallDir "$PROGRAMFILES64\\WinSnitch"
RequestExecutionLevel admin
Page directory
Page instfiles
Section "install"
  SetOutPath "$INSTDIR"
  File /r "..\\dist\\*"
  ExecWait '"$INSTDIR\\winsnitchd.exe" --install'
  ExecWait 'sc start WinSnitch'
SectionEnd
