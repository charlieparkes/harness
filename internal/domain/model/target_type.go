package model

//go:generate go tool go-enum

// TargetType is the kind of path a target points to.
// ENUM(Unspecified, Generic, VersionControlGit)
type TargetType string
