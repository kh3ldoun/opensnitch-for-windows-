//go:build windows

package wfp

import (
	"context"
	"time"

	"github.com/opensnitch/winsnitch/backend/internal/events"
	"github.com/sirupsen/logrus"
	"github.com/tailscale/wf"
)

type Engine struct {
	log *logrus.Entry
}

func New(log *logrus.Entry) *Engine {
	return &Engine{log: log}
}

func (e *Engine) Run(ctx context.Context, out chan<- events.ConnectionEvent) error {
	e.log.Info("starting WFP interception engine")
	sess, err := wf.New(&wf.Options{Name: "WinSnitch"})
	if err != nil {
		return err
	}
	defer sess.Close()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case t := <-ticker.C:
			out <- events.ConnectionEvent{
				ID:          t.Format(time.RFC3339Nano),
				ProcessPath: `C:\\Windows\\System32\\curl.exe`,
				Domain:      "example.org",
				DstIP:       "93.184.216.34",
				DstPort:     443,
				Protocol:    "tcp",
				Timestamp:   t,
				State:       "pending",
			}
		}
	}
}
