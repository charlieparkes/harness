package model

import (
	"time"

	"github.com/google/uuid"
)

// AgentActivity is an immutable event for one agent run.
type AgentActivity struct {
	ID      uuid.UUID         `json:"id"`
	AgentID string            `json:"agent_id"`
	Type    AgentActivityType `json:"type"`

	// Log represents a structured log.
	Log []byte `json:"log"`

	CreatedAt time.Time `json:"created_at"`
}
