package model

import "time"

// Plan is the plan document for a task.
type Plan struct {
	ID        string
	TaskID    string
	Revision  int
	Status    PlanStatus
	Document  PlanDocument
	CreatedAt time.Time
	UpdatedAt *time.Time
}

// PlanDocument is the JSONB payload for all plan revisions.
type PlanDocument struct{}
