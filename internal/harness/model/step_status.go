package model

//go:generate go tool go-enum

// StepStatus is the lifecycle state of a step.
// ENUM(Unspecified, Pending, Ready, Running, PendingFeedback, Completed, Canceled, Failed)
type StepStatus string
