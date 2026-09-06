package model

import "time"

// Agent is one agent run for a task.
type Agent struct {
	ID     string    `json:"id"`
	TaskID string    `json:"task_id"`
	Role   AgentRole `json:"role"`
	Model  string    `json:"model"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}
