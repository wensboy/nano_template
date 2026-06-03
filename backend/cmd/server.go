package main

import (
	"os"
	"os/signal"
	"syscall"

	"example.com/nano_template/pkg/config"
	"example.com/nano_template/pkg/util"
	"go.uber.org/zap"
)

type Server interface {
	Start(*config.Config)
	StartBg(*config.Config)
	Stop()
}

type RunnerOption func(*ServerRunner)

type ServerRunner struct {
	cfg    *config.Config
	server Server
}

func NewServerRunner(opts ...RunnerOption) *ServerRunner {
	runner := &ServerRunner{}
	for _, opt := range opts {
		opt(runner)
	}
	return runner
}

func WithConfig(cfg *config.Config) RunnerOption {
	return func(runner *ServerRunner) {
		runner.cfg = cfg
	}
}

func WithServer(server Server) RunnerOption {
	return func(runner *ServerRunner) {
		runner.server = server
	}
}

func (r *ServerRunner) Run() {
	if r.cfg == nil {
		util.Error("[runner]", zap.String("exec-run", "config not found"))
		os.Exit(1)
	}
	r.server.StartBg(r.cfg)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	r.server.Stop()
}
