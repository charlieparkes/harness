package model

import "time"

// Phase is one workflow phase of a task.
type Phase struct {
	ID        string
	TaskID    string
	Type      PhaseType
	Status    PhaseStatus
	Error     *string
	CreatedAt time.Time
	UpdatedAt *time.Time
}
