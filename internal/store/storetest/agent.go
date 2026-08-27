package storetest

import (
	"context"
	"testing"

	"github.com/charlieparkes/go-testsize"
	"github.com/charlieparkes/harness/internal/harness"
)

// AgentStoreTests checks GetRunningAgentCount on a harness.Store.
// newStore must return a clean store; it is called once per subtest.
func AgentStoreTests(t *testing.T, newStore func(*testing.T) harness.Store) {
	t.Helper()
	t.Run("GetRunningAgentCount", func(t *testing.T) {
		t.Parallel()

		t.Run("clean store is zero", func(t *testing.T) {
			t.Parallel()
			ctx := testsize.Small(t)
			store := newStore(t)

			got, err := store.GetRunningAgentCount(ctx)
			if err != nil {
				t.Fatalf("GetRunningAgentCount() error = %v", err)
			}
			if got != 0 {
				t.Fatalf("GetRunningAgentCount() = %d, want 0", got)
			}
		})

		t.Run("canceled context", func(t *testing.T) {
			t.Parallel()
			store := newStore(t)
			canceled := canceledContext(t)

			got, err := store.GetRunningAgentCount(canceled)
			requireErrIs(t, err, context.Canceled)
			if got != 0 {
				t.Fatalf("GetRunningAgentCount() = %d, want 0", got)
			}
		})
	})
}
