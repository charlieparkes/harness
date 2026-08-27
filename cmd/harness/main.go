package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/charlieparkes/harness/internal/harness/fake"
	"github.com/charlieparkes/harness/internal/harness/ui"
	"github.com/charlieparkes/harness/internal/logger"
	"github.com/charlieparkes/harness/internal/store/stub"
	"github.com/urfave/cli/v3"
)

func main() {
	l := logger.New()
	if err := run(l); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(logger *slog.Logger) (err error) {
	ctx, cancel := signal.NotifyContext(context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGHUP,
		syscall.SIGQUIT,
	)
	defer cancel()

	defer func() {
		if r := recover(); r != nil {
			cancel()
			logger.ErrorContext(ctx, "panic", "err", r)
		}
	}()

	store := stub.New()
	tasks := fake.Tasks(20)
	for _, t := range tasks {
		_, err := store.CreateTask(ctx, t)
		if err != nil {
			panic(err)
		}
	}

	cmd := &cli.Command{
		Name:  "harness",
		Usage: "create a consistent, repeatable experience when using agents to augment any task",
		Action: func(ctx context.Context, _ *cli.Command) error {
			return ui.Run(ctx, store)
		},
	}

	return cmd.Run(ctx, os.Args)
}
