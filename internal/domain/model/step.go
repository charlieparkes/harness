package model

import "time"

// Step is an executable unit from an accepted plan revision.
type Step struct {
	// ID should be initially set by the agent as something arbitrary
	// like "S1", "S2", etc. When a plan is accepted, the steps are given
	// real identifiers and written to the database for execution tracking.
	ID           string `json:"id"`
	PlanID       string `json:"plan_id"`
	PlanRevision int64  `json:"plan_revision"`

	// Step is not revisioned, and status is directly modified.
	// The expectations and description are not changed.
	Status          StepStatus       `json:"status"`
	Title           string           `json:"title"`
	Summary         string           `json:"summary"`
	ProposedChanges []ProposedChange `json:"proposed_changes"`
	ProposedTests   []ProposedTest   `json:"proposed_tests"`

	// list of verification steps taken and commands run
	Verifications []string `json:"verifications"`
	Risks         []string `json:"risks"`

	Dependencies []string `json:"depends_on"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}
