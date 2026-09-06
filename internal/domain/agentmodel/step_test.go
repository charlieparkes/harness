package agentmodel

import (
	"testing"

	"github.com/charlieparkes/go-testsize"
	"github.com/charlieparkes/harness/internal/domain/model"
	"github.com/charlieparkes/harness/test"
)

func TestStepStructFields(t *testing.T) {
	t.Parallel()
	testsize.Small(t)

	test.AssertStructFieldsEqual(t, Step{}, model.Step{},
		"CreatedAt",
		"UpdatedAt",
		"DeletedAt",
		"PlanID",
		"PlanRevision",
		"Status",
	)
}
