package model

//go:generate go tool go-enum

// PlanStatus is the lifecycle state of a plan.
// ENUM(Unspecified, Draft, Accepted)
type PlanStatus string
