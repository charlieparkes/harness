package agentmodel

import (
	"testing"

	"github.com/charlieparkes/go-testsize"
	"github.com/charlieparkes/harness/internal/domain/model"
	"github.com/charlieparkes/harness/test"
)

func TestReviewStructFields(t *testing.T) {
	t.Parallel()
	testsize.Small(t)

	test.AssertStructFieldsEqual(t, Review{}, model.Review{},
		"CreatedAt",
		"UpdatedAt",
		"DeletedAt",
		"ID",
		"Revision",
		"TaskID",
	)
}
