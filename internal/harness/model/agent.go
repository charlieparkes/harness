package model

import "time"

// Agent is one agent run for a task.
type Agent struct {
	ID        string
	TaskID    string
	Role      AgentRole
	Model     string
	CreatedAt time.Time
	UpdatedAt *time.Time
}
