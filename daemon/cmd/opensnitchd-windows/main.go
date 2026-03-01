//go:build windows

package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/evilsocket/opensnitch/daemon/log"
	"github.com/evilsocket/opensnitch/daemon/platform/windows/procinfo"
	"github.com/evilsocket/opensnitch/daemon/platform/windows/service"
	"github.com/evilsocket/opensnitch/daemon/platform/windows/wfp"
)

var (
	runAsService = true
	serviceName  = "OpenSnitch-Windows"
	configPath   = `C:\ProgramData\OpenSnitch-Windows\config\default-config.json`
)

func init() {
	flag.BoolVar(&runAsService, "service", runAsService, "Run as a Windows service")
	flag.StringVar(&serviceName, "service-name", serviceName, "Windows service name")
	flag.StringVar(&configPath, "config-file", configPath, "Path to daemon configuration")
}

func main() {
	flag.Parse()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	if err := log.OpenFile(log.StdoutFile); err != nil {
		panic(err)
	}
	defer log.Close()

	pi := procinfo.NewResolver()
	interceptor, err := wfp.NewInterceptor(wfp.Config{ConfigPath: configPath, ProcessResolver: pi})
	if err != nil {
		log.Fatal("failed to initialize WFP interceptor: %v", err)
	}
	defer interceptor.Close()

	daemon := service.NewDaemon(service.Config{
		Name:        serviceName,
		Interceptor: interceptor,
	})

	if runAsService {
		if err := daemon.RunService(ctx); err != nil {
			log.Fatal("service failed: %v", err)
		}
		return
	}

	if err := daemon.RunConsole(ctx); err != nil {
		log.Fatal("daemon failed: %v", err)
	}
}
