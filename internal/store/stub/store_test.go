package stub

import (
	"testing"

	"github.com/charlieparkes/go-testsize"
	"github.com/charlieparkes/harness/internal/domain"
	"github.com/charlieparkes/harness/internal/store/storetest"
)

func TestStore(t *testing.T) {
	t.Parallel()
	testsize.Small(t)
	storetest.StoreTests(t, func(_ *testing.T) domain.Store {
		return New()
	})
}
