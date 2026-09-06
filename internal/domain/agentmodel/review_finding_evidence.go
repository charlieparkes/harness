package agentmodel

import (
	"github.com/charlieparkes/harness/internal/domain/model"
)

// ReviewFindingEvidence describes one logical piece of evidence for a finding.
// Represents a "where", "when", "why", or "how."
type ReviewFindingEvidence struct {
	Description string       `json:"description"`
	Files       []model.File `json:"files"`
}

func NewReviewFindingEvidence(e model.ReviewFindingEvidence) ReviewFindingEvidence {
	files, _ := e.Files.Value()
	return ReviewFindingEvidence{
		Description: e.Description,
		Files:       files,
	}
}
