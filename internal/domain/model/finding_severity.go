package model

//go:generate go tool go-enum

// ReviewFindingSeverity is the severity of a review finding.
// ENUM(Unspecified, Low, Medium, High)
type ReviewFindingSeverity string
