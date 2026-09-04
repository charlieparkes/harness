package stub

import (
	"github.com/charlieparkes/harness/internal/harness/model"
)

// Store is a temporary in-memory store that returns generated stub tasks.
type Store struct {
	tasks []model.Task
}

func New() *Store {
	return &Store{tasks: make([]model.Task, 0)}
}
