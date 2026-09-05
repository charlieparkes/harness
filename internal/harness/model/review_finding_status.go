package model

//go:generate go tool go-enum

// ReviewFindingStatus is the disposition of a review finding.
// ENUM(Unspecified, Unresolved, Ignored, Rejected, Resolved)
type ReviewFindingStatus string
