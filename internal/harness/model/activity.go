package model

import (
	"time"

	"github.com/google/uuid"
)

// Activity is a timestamped event associated with a task.
type Activity struct {
	ID        uuid.UUID
	TaskID    string
	Type      ActivityType
	CreatedAt time.Time
}
