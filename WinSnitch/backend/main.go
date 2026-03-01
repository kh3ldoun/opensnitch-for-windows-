package main

import (
	"log"
	"os"

	"github.com/kardianos/service"
)

type program struct{}

func (p *program) Start(s service.Service) error {
	go p.run()
	return nil
}

func (p *program) run() {
	// Initialize WinSnitch Core
	log.Println("Starting WinSnitch backend daemon...")

	// Initialize WFP interceptor
	_, err := InitWFPInterceptor()
	if err != nil {
		log.Printf("Failed to init WFP: %v\n", err)
		return
	}

	// Wait for connections and communicate via RPC/WebSocket
}

func (p *program) Stop(s service.Service) error {
	log.Println("Stopping WinSnitch backend daemon...")
	return nil
}

func main() {
	svcConfig := &service.Config{
		Name:        "WinSnitch",
		DisplayName: "WinSnitch Backend Service",
		Description: "Interactive firewall daemon for WinSnitch.",
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}

	if len(os.Args) > 1 {
		err = service.Control(s, os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	err = s.Run()
	if err != nil {
		log.Fatal(err)
	}
}
