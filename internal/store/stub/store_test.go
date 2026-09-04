package stub

import (
	"testing"

	"github.com/charlieparkes/go-testsize"
	"github.com/charlieparkes/harness/internal/harness"
	"github.com/charlieparkes/harness/internal/store/storetest"
)

func TestStore(t *testing.T) {
	t.Parallel()
	testsize.Small(t)
	storetest.StoreTests(t, func(_ *testing.T) harness.Store {
		return New()
	})
}
