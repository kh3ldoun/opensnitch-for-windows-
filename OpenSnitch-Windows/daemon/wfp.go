package main

import (
	"log"
	"fmt"

	"github.com/tailscale/wf"
	"golang.org/x/sys/windows"
)

// WFPInterceptor struct wraps our access to WFP
type WFPInterceptor struct {
	session *wf.Session
}

// InitWFPInterceptor creates and configures a Windows Filtering Platform session for OpenSnitch.
// It returns a WFPInterceptor backed by a dynamic WFP session and configures initial filters.
// On error any partially created session is closed and the error is returned.
func InitWFPInterceptor() (*WFPInterceptor, error) {
	log.Println("Initializing Windows Filtering Platform...")

	// Open WFP session. Note we do dynamic so rules clean up on exit.
	session, err := wf.New(&wf.Options{
		Name:        "OpenSnitch Windows Daemon",
		Description: "Dynamically added rules for OpenSnitch connection filtering",
		Dynamic:     true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create WFP session: %w", err)
	}

	wfp := &WFPInterceptor{
		session: session,
	}

	err = wfp.setupFilters()
	if err != nil {
		session.Close()
		return nil, fmt.Errorf("failed to setup initial WFP filters: %w", err)
	}

	return wfp, nil
}

func (w *WFPInterceptor) setupFilters() error {
	// 1. Add our Sublayer
	err := w.session.AddSublayer(&wf.Sublayer{
		Name:        "OpenSnitch Block/Allow Sublayer",
		Description: "Sublayer for OpenSnitch decisions",
		Weight:      10000,
	})
	if err != nil {
		return fmt.Errorf("AddSublayer failed: %w", err)
	}
	log.Println("Added WFP sublayer.")

	// 2. We can add a simple default block rule for testing, but in a real scenario
	// we want a WFP callout driver. Since we are porting OpenSnitch, we assume
	// there's a kernel driver named "OpenSnitchCallout" that registers Callout IDs.
	// For now, we will just establish the framework.

	// Example of adding a basic rule using tailscale/wf
	// Here we would normally bind to the callout ID registered by our driver.

	return nil
}

func (w *WFPInterceptor) Close() {
	if w.session != nil {
		w.session.Close()
		log.Println("Closed WFP session.")
	}
}

// Map of common layer IDs
var (
	FWPM_LAYER_ALE_AUTH_CONNECT_V4 = windows.GUID{Data1: 0xc38d57d1, Data2: 0x05a7, Data3: 0x4c33, Data4: [8]byte{0x90, 0x4f, 0x7f, 0xbc, 0xee, 0xe6, 0x0e, 0x82}}
	FWPM_LAYER_ALE_AUTH_CONNECT_V6 = windows.GUID{Data1: 0x4a7239ce, Data2: 0xdce7, Data3: 0x4a36, Data4: [8]byte{0x81, 0x89, 0x06, 0x99, 0xc9, 0x74, 0xe2, 0x20}}
)
