package model

import "time"

// Review is one review pass for a task.
type Review struct {
	ID        string
	TaskID    string
	Pass      int
	Revision  int
	Document  ReviewDocument
	Decision  ReviewDecision
	CreatedAt time.Time
	UpdatedAt *time.Time
}

// ReviewDocument is the JSONB payload for all revisions of a review.
type ReviewDocument struct{}
