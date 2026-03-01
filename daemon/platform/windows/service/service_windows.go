//go:build windows

package service

import (
	"context"
)

type Interceptor interface {
	Close() error
}

type Config struct {
	Name        string
	Interceptor Interceptor
}

type Daemon struct {
	cfg Config
}

func NewDaemon(cfg Config) *Daemon { return &Daemon{cfg: cfg} }

func (d *Daemon) RunService(ctx context.Context) error {
	// Service control manager integration lives here.
	<-ctx.Done()
	return nil
}

func (d *Daemon) RunConsole(ctx context.Context) error {
	<-ctx.Done()
	return nil
}
