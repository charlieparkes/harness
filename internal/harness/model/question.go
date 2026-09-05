package model

// Question is a prompt that waits for a user answer.
type Question struct {
	ID               string
	Question         string
	SuggestedAnswers SuggestedAnswers
	Answer           *string
}

// SuggestedAnswers is a list of answer options for a question.
type SuggestedAnswers struct {
	Answers []string
}
