package main

import (
	"log"
	"sync"
	"time"

	pb "github.com/evilsocket/opensnitch-windows/daemon/proto"
)

// DecisionRequest represents a pending connection from the kernel/WFP.
type DecisionRequest struct {
	ID        uint64
	Pid       uint32
	Protocol  uint32
	SrcIP     string
	DstIP     string
	SrcPort   uint16
	DstPort   uint16
	Path      string
	Cmdline   string
	UserSID   string
	ResultCh  chan string // Action string back to kernel
}

// DecisionEngine manages asking the UI and returning decisions.
type DecisionEngine struct {
	pending map[uint64]*DecisionRequest
	mu      sync.Mutex
	uiConn  pb.UIServer // Interface to call UI, in real world we use streams
}

func NewDecisionEngine() *DecisionEngine {
	return &DecisionEngine{
		pending: make(map[uint64]*DecisionRequest),
	}
}

// HandleNewConnection is called by the WFP interceptor or Driver bridge.
func (de *DecisionEngine) HandleNewConnection(req *DecisionRequest) {
	de.mu.Lock()
	de.pending[req.ID] = req
	de.mu.Unlock()

	// 1. Gather process info
	path, _ := getProcessPath(req.Pid)
	user, _ := getProcessUser(req.Pid)
	cmd, _  := getCommandLine(req.Pid)

	req.Path = path
	req.UserSID = user
	req.Cmdline = cmd

	// 3. Ask UI
	log.Printf("Pending connection: PID %d (%s) -> %s:%d", req.Pid, req.Path, req.DstIP, req.DstPort)

	// Simulated timeout and default action:
	go func() {
		time.Sleep(15 * time.Second)
		de.mu.Lock()
		if _, exists := de.pending[req.ID]; exists {
			delete(de.pending, req.ID)
			log.Printf("Connection %d timed out. Default action: DENY", req.ID)
			req.ResultCh <- "deny"
		}
		de.mu.Unlock()
	}()
}

// UIDecision is called when the Python UI sends a reply to a prompt.
func (de *DecisionEngine) UIDecision(id uint64, action string, duration string) {
	de.mu.Lock()
	defer de.mu.Unlock()

	if req, exists := de.pending[id]; exists {
		log.Printf("UI replied for connection %d: Action=%v, Duration=%s", id, action, duration)
		req.ResultCh <- action
		delete(de.pending, id)
	}
}
