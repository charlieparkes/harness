package ui

import (
	"context"
	"errors"

	tea "charm.land/bubbletea/v2"
)

func Run(ctx context.Context, store Store) error {
	program := tea.NewProgram(New(ctx, store), tea.WithContext(ctx))
	_, err := program.Run()
	if err == nil {
		return nil
	}
	if errors.Is(err, tea.ErrInterrupted) || errors.Is(err, tea.ErrProgramKilled) {
		return nil
	}
	return err
}
