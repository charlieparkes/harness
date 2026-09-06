package model

//go:generate go tool go-enum

// ReviewFindingCategory is the disposition of a review finding.
// ENUM(Unspecified, Bug)
type ReviewFindingCategory string
