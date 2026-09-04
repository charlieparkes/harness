package storetest

import (
	"testing"

	"github.com/charlieparkes/harness/internal/harness"
)

func StoreTests(t *testing.T, factory func(*testing.T) harness.Store) {
	t.Helper()
	t.Run("Task", func(t *testing.T) {
		t.Parallel()
		TaskStoreTests(t, factory)
	})
	t.Run("Agent", func(t *testing.T) {
		t.Parallel()
		AgentStoreTests(t, factory)
	})
}
