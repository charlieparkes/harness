package model

import "time"

// Phase is one workflow phase of a task.
type Phase struct {
	ID     string      `json:"id"`
	TaskID string      `json:"task_id"`
	Type   PhaseType   `json:"type"`
	Status PhaseStatus `json:"status"`
	Error  *string     `json:"error"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}
