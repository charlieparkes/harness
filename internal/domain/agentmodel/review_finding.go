package agentmodel

import (
	"github.com/charlieparkes/go-transform"
	"github.com/charlieparkes/harness/internal/domain/model"
)

// ReviewFinding is a finding recorded against a review.
type ReviewFinding struct {
	// Agent-sourced identifier for this review finding
	// For example, R1, R2, R3, etc
	ID string `json:"id"`

	Severity        model.ReviewFindingSeverity `json:"severity"`
	Category        model.ReviewFindingCategory `json:"category"`
	Evidence        []ReviewFindingEvidence     `json:"evidence"`
	ProposedChanges []ProposedChange            `json:"proposed_changes"`
	ProposedTests   []ProposedTest              `json:"proposed_tests"`
}

func NewReviewFinding(f model.ReviewFinding) ReviewFinding {
	evidences, _ := f.Evidence.Value()
	changes, _ := f.ProposedChanges.Value()
	tests, _ := f.ProposedTests.Value()
	return ReviewFinding{
		ID:              f.ID,
		Severity:        f.Severity,
		Category:        f.Category,
		Evidence:        transform.Slice(evidences, NewReviewFindingEvidence),
		ProposedChanges: transform.Slice(changes, NewProposedChange),
		ProposedTests:   transform.Slice(tests, NewProposedTest),
	}
}
