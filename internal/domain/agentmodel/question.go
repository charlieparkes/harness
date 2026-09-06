package agentmodel

import (
	"github.com/charlieparkes/go-transform"
	"github.com/charlieparkes/harness/internal/domain/model"
)

// Questions is a prompt that waits for a user answers to multiple questions.
type Questions struct {
	Questions []Question `json:"questions"`
}

// Question is a prompt that waits for a user answer.
type Question struct {
	// Agent-sourced identifier for this question
	// For example, Q1, Q2, Q3, etc
	ID string `json:"id"`

	// One to three sentence question
	Question string `json:"question"`

	// Agent-sourced suggested answers
	SuggestedAnswers []string `json:"suggested_answers"`
}

func NewQuestion(q model.Question) Question {
	return Question{
		ID:               q.ID,
		Question:         q.Question,
		SuggestedAnswers: q.SuggestedAnswers,
	}
}

func NewQuestions(q []model.Question) Questions {
	return Questions{Questions: transform.Slice(q, NewQuestion)}
}
