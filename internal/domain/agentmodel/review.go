package agentmodel

import (
	"github.com/charlieparkes/go-transform"
	"github.com/charlieparkes/harness/internal/domain/model"
)

// Review is a multi-pass review of a task.
// Each pass of review will result in creating one or more [ReviewFinding].
type Review struct {
	Decision model.ReviewDecision `json:"decision"`

	// list of non-blocking observations made
	Observations []string `json:"observations"`

	// list of verification steps taken and commands run
	Verifications []string `json:"verifications"`

	Findings []ReviewFinding `json:"findings"`
}

func NewReview(r model.Review) Review {
	decision, _ := r.Decision.Value()
	observations, _ := r.Observations.Value()
	verifications, _ := r.Verifications.Value()
	findings, _ := r.Findings.Value()
	return Review{
		Decision:      decision,
		Observations:  observations,
		Verifications: verifications,
		Findings:      transform.Slice(findings, NewReviewFinding),
	}
}
