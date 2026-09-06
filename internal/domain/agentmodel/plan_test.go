package agentmodel

import (
	"testing"

	"github.com/charlieparkes/go-testsize"
	"github.com/charlieparkes/harness/internal/domain/model"
	"github.com/charlieparkes/harness/test"
)

func TestPlanStructFields(t *testing.T) {
	t.Parallel()
	testsize.Small(t)

	test.AssertStructFieldsEqual(t, Plan{}, model.Plan{},
		"CreatedAt",
		"UpdatedAt",
		"DeletedAt",
		"ID",
		"Revision",
		"TaskID",
		"Status",
	)
}
