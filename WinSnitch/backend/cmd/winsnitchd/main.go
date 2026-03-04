package main

import (
	"flag"
	"fmt"
	"os"

	kservice "github.com/kardianos/service"
	wsservice "github.com/opensnitch/winsnitch/backend/internal/service"
)

func main() {
	install := flag.Bool("install", false, "install windows service")
	uninstall := flag.Bool("uninstall", false, "remove windows service")
	console := flag.Bool("console", false, "run in foreground")
	flag.Parse()

	program := wsservice.NewProgram()
	svcCfg := &kservice.Config{
		Name:        "WinSnitch",
		DisplayName: "WinSnitch Firewall",
		Description: "Interactive outbound firewall for Windows",
	}
	svc, err := kservice.New(program, svcCfg)
	if err != nil {
		panic(err)
	}

	if *install {
		must(svc.Install())
		fmt.Println("WinSnitch service installed")
		return
	}
	if *uninstall {
		must(svc.Uninstall())
		fmt.Println("WinSnitch service uninstalled")
		return
	}
	if *console {
		must(wsservice.RunInteractive())
		return
	}
	must(svc.Run())
}

func must(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
