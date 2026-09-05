package model

//go:generate go tool go-enum

// WorktreeStatus is the lifecycle state of a worktree.
// ENUM(Unspecified, Pending, Preparing, Ready, Removed)
type WorktreeStatus string
