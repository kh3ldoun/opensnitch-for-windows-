package service

import (
	"context"
	kservice "github.com/kardianos/service"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/opensnitch/winsnitch/backend/internal/api"
	"github.com/opensnitch/winsnitch/backend/internal/blocklist"
	"github.com/opensnitch/winsnitch/backend/internal/config"
	"github.com/opensnitch/winsnitch/backend/internal/events"
	"github.com/opensnitch/winsnitch/backend/internal/logging"
	"github.com/opensnitch/winsnitch/backend/internal/rules"
	"github.com/opensnitch/winsnitch/backend/internal/wfp"
	"github.com/sirupsen/logrus"
)

type Program struct {
	cfg    config.Settings
	log    *logrus.Logger
	cancel context.CancelFunc
	done   chan struct{}
}

func NewProgram() *Program {
	return &Program{done: make(chan struct{})}
}

func (p *Program) Start(_ kservice.Service) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(cfg.ProgramDataDir, 0o755); err != nil {
		return err
	}
	log, err := logging.New(cfg.ProgramDataDir)
	if err != nil {
		return err
	}
	p.cfg, p.log = cfg, log

	ctx, cancel := context.WithCancel(context.Background())
	p.cancel = cancel
	go p.run(ctx)
	return nil
}

func (p *Program) run(ctx context.Context) {
	defer close(p.done)
	entry := logrus.NewEntry(p.log)
	rulesStore, err := rules.NewStore(p.cfg.ProgramDataDir)
	if err != nil {
		entry.WithError(err).Error("rules init failed")
		return
	}
	_ = rulesStore

	bl := blocklist.New()
	if err := bl.Update(p.cfg.BlocklistsURL); err != nil {
		entry.WithError(err).Warn("blocklist initial update failed")
	}

	hub := api.NewHub(entry)
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", hub.Handler)
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) { _, _ = w.Write([]byte("ok")) })

	srv := &http.Server{Addr: p.cfg.ListenAddress, Handler: mux}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			entry.WithError(err).Error("api server failed")
		}
	}()

	e := wfp.New(entry)
	eventsCh := make(chan events.ConnectionEvent, 128)
	go func() {
		if err := e.Run(ctx, eventsCh); err != nil {
			entry.WithError(err).Warn("interception engine exited")
		}
	}()

	for {
		select {
		case <-ctx.Done():
			_ = srv.Shutdown(context.Background())
			hub.Close(context.Background())
			return
		case ev := <-eventsCh:
			if bl.Contains(ev.Domain) {
				ev.State = "blocked_by_blocklist"
			}
			hub.Broadcast(ev)
		}
	}
}

func (p *Program) Stop(_ kservice.Service) error {
	if p.cancel != nil {
		p.cancel()
	}
	select {
	case <-p.done:
	case <-time.After(10 * time.Second):
	}
	return nil
}

func RunInteractive() error {
	prog := NewProgram()
	if err := prog.Start(kservice.Service(nil)); err != nil {
		return err
	}
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	return prog.Stop(kservice.Service(nil))
}

func ConfigPath(programDataDir string) string {
	return filepath.Join(programDataDir, "winsnitch.yaml")
}
