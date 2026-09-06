package model

//go:generate go tool go-enum

// PhaseType is the kind of workflow phase.
// ENUM(Unspecified, Worktree, Plan, Execute, Review, Apply)
type PhaseType string
