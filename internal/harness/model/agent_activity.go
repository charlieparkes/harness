package model

import "time"

// AgentActivity is an immutable event for one agent run.
type AgentActivity struct {
	ID        string
	AgentID   string
	Type      AgentActivityType
	Log       AgentActivityLog
	CreatedAt time.Time
}

// AgentActivityLog is the JSONB payload for a structured event log.
type AgentActivityLog struct{}
