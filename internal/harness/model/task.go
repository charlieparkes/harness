package model

import (
	"time"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

// Task is the minimal identity and display text shown in the task list.
type Task struct {
	ID        string
	Title     string
	Status    TaskStatus
	CreatedAt time.Time
	UpdatedAt *time.Time
}

// NewTask returns a ready task with a Crockford Base32 nanoid ID and the current UTC time.
func NewTask(title string) (Task, error) {
	id, err := gonanoid.Generate(gonanoid.CrockfordBase32Upper, 4)
	if err != nil {
		return Task{}, err
	}
	return Task{
		ID:        id,
		Title:     title,
		Status:    TaskStatusReady,
		CreatedAt: time.Now().UTC(),
	}, nil
}

// LastActive returns UpdatedAt when set, otherwise CreatedAt.
func (t Task) LastActive() time.Time {
	if t.UpdatedAt != nil {
		return *t.UpdatedAt
	}
	return t.CreatedAt
}

// IsActive reports whether the task is in an active status and has not been canceled.
func (t Task) IsActive() bool {
	switch t.Status {
	case TaskStatusReady,
		TaskStatusRunning,
		TaskStatusPendingFeedback,
		TaskStatusFailed:
		return true
	case TaskStatusUnspecified,
		TaskStatusCompleted,
		TaskStatusCanceled:
		return false
	default:
		return false
	}
}

// Compare orders active tasks before inactive tasks, then later LastActive values first.
func (t Task) Compare(t2 Task) int {
	a := t.IsActive()
	a2 := t2.IsActive()
	switch {
	case a && !a2:
		return -1
	case !a && a2:
		return 1
	}
	return t2.LastActive().Compare(t.LastActive())
}
