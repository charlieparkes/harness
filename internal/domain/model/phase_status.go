package model

//go:generate go tool go-enum

// PhaseStatus is the lifecycle state of a phase.
// ENUM(Unspecified, Pending, Ready, Running, Feedback, Completed, Canceled, Failed)
type PhaseStatus string
