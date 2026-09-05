package model

import "time"

// Step is an executable unit from an accepted plan revision.
type Step struct {
	ID           string
	PlanID       string
	PlanRevision int
	Title        string
	Description  string
	Status       StepStatus
	Result       StepResult
	CreatedAt    time.Time
	UpdatedAt    *time.Time
}

// StepResult is the JSONB payload for a structured agent result.
type StepResult struct{}
