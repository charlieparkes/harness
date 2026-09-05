package model

//go:generate go tool go-enum

// AgentRole is the role of an agent run.
// ENUM(Unspecified, Planner, PlanReviewer, Executor, Reviewer)
type AgentRole string
