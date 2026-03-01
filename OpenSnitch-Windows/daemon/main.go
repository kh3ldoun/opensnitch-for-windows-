package main

import (
	"context"
	"flag"

	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/debug"

)

var (
	logFile   = flag.String("log-file", filepath.Join(os.Getenv("ProgramData"), "OpenSnitch-Windows", "daemon.log"), "Path to log file")
	rulesPath = flag.String("rules-path", filepath.Join(os.Getenv("ProgramData"), "OpenSnitch-Windows", "rules"), "Path to rules directory")
	logLevel  = flag.Int("log-level", 2, "Log level (0=quiet, 1=error, 2=warn, 3=info, 4=debug)")
)

const svcName = "OpenSnitch"

// winService implements the svc.Handler interface
type winService struct{}

func (m *winService) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {
	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown
	changes <- svc.Status{State: svc.StartPending}

	// Initialize OpenSnitch Core
	log.Printf("Starting OpenSnitch daemon...")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize WFP Interceptor
	interceptor, err := InitWFPInterceptor()
	if err != nil {
		log.Printf("Failed to initialize WFP interceptor: %v", err)
		return
	}
	defer interceptor.Close()

	// Start UI gRPC server
	go startGRPCServer(ctx)

	changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}

	for {
		select {
		case c := <-r:
			switch c.Cmd {
			case svc.Interrogate:
				changes <- c.CurrentStatus
			case svc.Stop, svc.Shutdown:
				log.Println("Received stop signal. Shutting down...")
				changes <- svc.Status{State: svc.StopPending}
				cancel()
				return
			default:
				log.Printf("Unexpected control request #%d", c)
			}
		case <-ctx.Done():
			return
		}
	}
}

func runService(name string, isDebug bool) {
	var err error
	if isDebug {
		err = debug.Run(name, &winService{})
	} else {
		err = svc.Run(name, &winService{})
	}
	if err != nil {
		log.Fatalf("Service %s failed: %v", name, err)
	}
}

func main() {
	flag.Parse()

	// Setup logging
	os.MkdirAll(filepath.Dir(*logFile), 0755)
	f, err := os.OpenFile(*logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	inService, err := svc.IsWindowsService()
	if err != nil {
		log.Fatalf("failed to determine if we are running in an interactive session: %v", err)
	}

	if inService {
		runService(svcName, false)
		return
	}

	// Interactive/console mode
	log.Println("Running in console mode. Press Ctrl+C to exit.")
	ctx, cancel := context.WithCancel(context.Background())

	// Setup signal handling
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	// Initialize WFP
	interceptor, err := InitWFPInterceptor()
	if err != nil {
		log.Fatalf("Failed to initialize WFP interceptor: %v", err)
	}
	defer interceptor.Close()

	go startGRPCServer(ctx)

	<-sigCh
	log.Println("Shutting down...")
	cancel()
	time.Sleep(1 * time.Second)
}
