package model

import "time"

// Target is a file or directory path associated with a task.
type Target struct {
	ID        string
	TaskID    string
	Path      string
	Type      TargetType
	Purpose   TargetPurpose
	CreatedAt time.Time
	UpdatedAt *time.Time
}
