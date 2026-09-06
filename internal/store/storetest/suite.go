package storetest

import (
	"testing"

	"github.com/charlieparkes/harness/internal/domain"
)

func StoreTests(t *testing.T, factory func(*testing.T) domain.Store) {
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
