package agentmodel

import (
	"testing"

	"github.com/charlieparkes/go-testsize"
	"github.com/charlieparkes/harness/internal/domain/model"
	"github.com/charlieparkes/harness/test"
)

func TestProposedChangeStructFields(t *testing.T) {
	t.Parallel()
	testsize.Small(t)

	test.AssertStructFieldsEqual(t, ProposedChange{}, model.ProposedChange{})
}

func TestProposedTestStructFields(t *testing.T) {
	t.Parallel()
	testsize.Small(t)

	test.AssertStructFieldsEqual(t, ProposedTest{}, model.ProposedTest{})
}
