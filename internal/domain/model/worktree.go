package model

import "time"

// Worktree is a Git worktree created for a VersionControlGit target.
type Worktree struct {
	ID           string         `json:"id"`
	TargetID     string         `json:"target_id"`
	StartRef     string         `json:"start_ref"`
	HeadCommit   *string        `json:"head_commit"`
	WorktreePath *string        `json:"worktree_path"`
	Status       WorktreeStatus `json:"status"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}
