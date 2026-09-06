package model

import (
	"time"

	"github.com/charlieparkes/harness/internal/lineage"
)

// ReviewFinding is a finding recorded against a review.
type ReviewFinding struct {
	ID       string `json:"id"`
	Revision int64  `json:"revision"`
	ReviewID string `json:"review_id"`

	Severity        ReviewFindingSeverity                            `json:"severity"`
	Category        ReviewFindingCategory                            `json:"category"`
	Status          lineage.RevisionedField[ReviewFindingStatus]     `json:"status"`
	Evidence        lineage.RevisionedField[[]ReviewFindingEvidence] `json:"evidence"`
	ProposedChanges lineage.RevisionedField[[]ProposedChange]        `json:"proposed_changes"`
	ProposedTests   lineage.RevisionedField[[]ProposedTest]          `json:"proposed_tests"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

// ReviewFindingEvidence describes one logical piece of evidence for a finding.
// Represents a "where", "when", "why", or "how."
type ReviewFindingEvidence struct {
	Description string                          `json:"description"`
	Files       lineage.RevisionedField[[]File] `json:"files"`
}
