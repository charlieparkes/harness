package model

// ReviewFinding is a finding recorded against a review.
type ReviewFinding struct {
	ID       string
	ReviewID string
	Severity FindingSeverity
	Status   ReviewFindingStatus
	Document ReviewFindingDocument
}

// ReviewFindingDocument is the JSONB payload for a finding snapshot.
type ReviewFindingDocument struct{}
