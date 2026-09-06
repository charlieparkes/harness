package model

import "time"

// Question is a prompt that waits for a user answer.
type Question struct {
	ID string `json:"id"`

	Question         string   `json:"question"`
	SuggestedAnswers []string `json:"suggested_answers"`
	Answer           *string  `json:"answer"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}
