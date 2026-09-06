package model

//go:generate go tool go-enum

// AutomationPolicyType is the approval mode for a phase.
// ENUM(Unspecified, Manual, Automatic)
type AutomationPolicyType string
