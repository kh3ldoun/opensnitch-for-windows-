package main

import (
	"fmt"
	"log"

	"github.com/tailscale/wf"
)

// WFPInterceptor holds the WFP Session
type WFPInterceptor struct {
	session *wf.Session
}

func InitWFPInterceptor() (*WFPInterceptor, error) {
	log.Println("Initializing Windows Filtering Platform (WFP)...")

	// Create dynamic WFP session (rules clean up on exit)
	session, err := wf.New(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create WFP session: %w", err)
	}

	wfpObj := &WFPInterceptor{
		session: session,
	}

	// Add a sublayer to hold our rules
	err = wfpObj.session.AddSublayer(&wf.Sublayer{
		Name:        "WinSnitch WFP Sublayer",
		Description: "Sublayer for WinSnitch decisions",
		Weight:      10000,
	})
	if err != nil {
		session.Close()
		return nil, fmt.Errorf("failed to add WFP sublayer: %w", err)
	}

	log.Println("Successfully created WFP session and sublayer.")

	return wfpObj, nil
}

func (w *WFPInterceptor) Close() {
	if w.session != nil {
		w.session.Close()
	}
}
