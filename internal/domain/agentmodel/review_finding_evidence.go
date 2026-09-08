package agentmodel

import (
	"github.com/charlieparkes/harness/internal/domain/model"
	"github.com/charlieparkes/harness/internal/lineage"
)

// ReviewFindingEvidence describes one logical piece of evidence for a finding.
// Represents a "where", "when", "why", or "how."
type ReviewFindingEvidence struct {
	Description string       `json:"description"`
	Files       []model.File `json:"files"`
}

func NewReviewFindingEvidence(e model.ReviewFindingEvidence) ReviewFindingEvidence {
	files := e.Files.Value()
	return ReviewFindingEvidence{
		Description: e.Description,
		Files:       files,
	}
}

func (e ReviewFindingEvidence) model() model.ReviewFindingEvidence {
	return model.ReviewFindingEvidence{
		Description: e.Description,
		Files:       lineage.NewRevisionedField(e.Files),
	}
}
