package model

import "time"

// Target is a file or directory path associated with a task.
type Target struct {
	ID      string        `json:"id"`
	TaskID  string        `json:"task_id"`
	Path    string        `json:"path"`
	Type    TargetType    `json:"type"`
	Purpose TargetPurpose `json:"purpose"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}
