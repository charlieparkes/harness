package model

import "time"

// Worktree is a Git worktree created for a VersionControlGit target.
type Worktree struct {
	ID           string
	TargetID     string
	StartRef     string
	HeadCommit   *string
	WorktreePath *string
	Status       WorktreeStatus
	CreatedAt    time.Time
	UpdatedAt    *time.Time
}
