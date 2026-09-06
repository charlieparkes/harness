package model

import (
	"time"

	"github.com/charlieparkes/harness/internal/lineage"
)

// Review is a multi-pass review of a task.
// Each pass of review will result in creating one or more [ReviewFinding].
type Review struct {
	ID       string `json:"id"`
	Revision int64  `json:"revision"` // new revision for each "pass"
	TaskID   string `json:"task_id"`

	Decision lineage.RevisionedField[ReviewDecision] `json:"decision"`

	// list of non-blocking observations made
	Observations lineage.RevisionedField[[]string] `json:"observations"`

	// list of verification steps taken and commands run
	Verifications lineage.RevisionedField[[]string] `json:"verifications"`

	Findings lineage.RevisionedField[[]ReviewFinding] `json:"findings"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}
