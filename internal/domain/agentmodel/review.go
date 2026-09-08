package agentmodel

import (
	"time"

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
	decision := r.Decision.Value()
	observations := r.Observations.Value()
	verifications := r.Verifications.Value()
	findings := r.Findings.Value()
	return Review{
		Decision:      decision,
		Observations:  observations,
		Verifications: verifications,
		Findings:      transform.Slice(findings, NewReviewFinding),
	}
}

func (r Review) Apply(prev model.Review, now time.Time) model.Review {
	rev := prev.Revision + 1

	changed := prev.Decision.Set(rev, r.Decision)
	changed = prev.Observations.Set(rev, r.Observations) || changed
	changed = prev.Verifications.Set(rev, r.Verifications) || changed
	changed = prev.Findings.Set(rev, transform.Slice(r.Findings, ReviewFinding.model)) || changed

	if changed {
		prev.Revision = rev
		prev.UpdatedAt = &now
	}
	return prev
}
