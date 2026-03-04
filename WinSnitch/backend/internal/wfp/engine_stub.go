//go:build !windows

package wfp

import (
	"context"
	"errors"

	"github.com/opensnitch/winsnitch/backend/internal/events"
	"github.com/sirupsen/logrus"
)

type Engine struct {
	log *logrus.Entry
}

func New(log *logrus.Entry) *Engine {
	return &Engine{log: log}
}

func (e *Engine) Run(ctx context.Context, out chan<- events.ConnectionEvent) error {
	_ = ctx
	_ = out
	return errors.New("WFP engine requires windows")
}
