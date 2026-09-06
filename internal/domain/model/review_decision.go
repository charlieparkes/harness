package model

//go:generate go tool go-enum

// ReviewDecision is the outcome of a review pass.
// ENUM(Unspecified, ChangesRequested, Approved, Rejected)
type ReviewDecision string
