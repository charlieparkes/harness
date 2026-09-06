package agentmodel

import (
	"testing"

	"github.com/charlieparkes/go-testsize"
	"github.com/charlieparkes/harness/internal/domain/model"
	"github.com/charlieparkes/harness/test"
)

func TestQuestionStructFields(t *testing.T) {
	t.Parallel()
	testsize.Small(t)

	test.AssertStructFieldsEqual(t, Question{}, model.Question{},
		"CreatedAt",
		"UpdatedAt",
		"DeletedAt",
		"Answer",
	)
}
