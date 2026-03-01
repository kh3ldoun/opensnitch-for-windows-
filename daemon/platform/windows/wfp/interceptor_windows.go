//go:build windows

package wfp

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/tailscale/wf"
)

type ProcessResolver interface {
	Resolve(pid uint32) (ProcessMetadata, error)
}

type ProcessMetadata struct {
	PID         uint32
	Path        string
	CommandLine string
	UserSID     string
}

type Config struct {
	ConfigPath      string
	ProcessResolver ProcessResolver
}

type Decision struct {
	Permit    bool
	Permanent bool
	RuleID    string
}

type Interceptor struct {
	cfg      Config
	session  *wf.Session
	provider wf.Provider
	subLayer wf.Sublayer

	mu      sync.RWMutex
	pending map[string]chan Decision
}

func NewInterceptor(cfg Config) (*Interceptor, error) {
	s, err := wf.New(&wf.Options{Name: "OpenSnitch-Windows"})
	if err != nil {
		return nil, fmt.Errorf("create wfp session: %w", err)
	}

	i := &Interceptor{cfg: cfg, session: s, pending: map[string]chan Decision{}}
	if err := i.bootstrap(); err != nil {
		s.Close()
		return nil, err
	}
	return i, nil
}

func (i *Interceptor) bootstrap() error {
	provider, err := i.session.AddProvider(&wf.Provider{Name: "OpenSnitch Provider", Description: "OpenSnitch-Windows policy provider"})
	if err != nil {
		return fmt.Errorf("add provider: %w", err)
	}
	sublayer, err := i.session.AddSublayer(&wf.Sublayer{Name: "OpenSnitch Sublayer", Provider: provider.ID(), Weight: 0x7fff})
	if err != nil {
		return fmt.Errorf("add sublayer: %w", err)
	}

	// Baseline block rules are inserted at ALE_AUTH_CONNECT layers.
	// Runtime decisions are expected to be fulfilled by the callout driver (driver/).
	_, err = i.session.AddFilter(&wf.Filter{
		Name:       "OpenSnitch Connect V4",
		Layer:      wf.LayerALEAuthConnectV4,
		Sublayer:   sublayer.ID(),
		Action:     wf.ActionPermit,
		Weight:     wf.EmptyWeight(),
		Provider:   provider.ID(),
		Persistent: true,
	})
	if err != nil {
		return fmt.Errorf("add ipv4 filter: %w", err)
	}

	_, err = i.session.AddFilter(&wf.Filter{
		Name:       "OpenSnitch Connect V6",
		Layer:      wf.LayerALEAuthConnectV6,
		Sublayer:   sublayer.ID(),
		Action:     wf.ActionPermit,
		Weight:     wf.EmptyWeight(),
		Provider:   provider.ID(),
		Persistent: true,
	})
	if err != nil {
		return fmt.Errorf("add ipv6 filter: %w", err)
	}

	i.provider = provider
	i.subLayer = sublayer
	return nil
}

func (i *Interceptor) Pend(requestID string) <-chan Decision {
	i.mu.Lock()
	defer i.mu.Unlock()
	ch := make(chan Decision, 1)
	i.pending[requestID] = ch
	return ch
}

func (i *Interceptor) Resolve(requestID string, decision Decision) error {
	i.mu.Lock()
	defer i.mu.Unlock()
	ch, ok := i.pending[requestID]
	if !ok {
		return errors.New("request not found")
	}
	ch <- decision
	close(ch)
	delete(i.pending, requestID)
	return nil
}

func (i *Interceptor) WaitDecision(ctx context.Context, requestID string, timeout time.Duration) (Decision, error) {
	ch := i.Pend(requestID)
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return Decision{}, ctx.Err()
	case <-timer.C:
		return Decision{}, errors.New("decision timeout")
	case d := <-ch:
		return d, nil
	}
}

func (i *Interceptor) Close() error {
	if i.session != nil {
		return i.session.Close()
	}
	return nil
}
