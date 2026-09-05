package model

//go:generate go tool go-enum

// FindingSeverity is the severity of a review finding.
// ENUM(Unspecified, High, Medium, Low)
type FindingSeverity string
