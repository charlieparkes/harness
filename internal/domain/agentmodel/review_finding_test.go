package agentmodel

import (
	"testing"

	"github.com/charlieparkes/go-testsize"
	"github.com/charlieparkes/harness/internal/domain/model"
	"github.com/charlieparkes/harness/test"
)

func TestReviewFindingStructFields(t *testing.T) {
	t.Parallel()
	testsize.Small(t)

	test.AssertStructFieldsEqual(t, ReviewFinding{}, model.ReviewFinding{},
		"CreatedAt",
		"UpdatedAt",
		"DeletedAt",
		"Revision",
		"ReviewID",
		"Status",
	)
}
