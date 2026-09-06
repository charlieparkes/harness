package model

// PlanQuestion links a question to a plan revision.
type PlanQuestion struct {
	PlanID       string
	PlanRevision int
	QuestionID   string
}
