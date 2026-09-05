package model

// StepDependency is a prerequisite edge between two steps.
type StepDependency struct {
	StepID          string
	DependsOnStepID string
}
