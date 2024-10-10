package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func Run(ctx context.Context, args []string) error {
	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	var cmd string
	if len(args) > 0 {
		cmd, args = args[0], args[1:]
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: false,
		Level:     slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	switch cmd {
	case "worker":
		return NewWorker().Run(ctx, args)
	default:
		return fmt.Errorf("%s %s: unknown command", Name, cmd)
	}
}
