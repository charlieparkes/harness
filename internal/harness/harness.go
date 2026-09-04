package harness

import (
	"context"
	"errors"

	"github.com/charlieparkes/harness/internal/harness/model"
)

// ErrTaskNotFound is returned when GetTask cannot find a task with the given ID.
var ErrTaskNotFound = errors.New("task not found")

// Store is the harness-layer persistence contract.
type Store interface {
	ListTasks(ctx context.Context) ([]model.Task, error)
	GetTask(ctx context.Context, id string) (model.Task, error)
	CreateTask(ctx context.Context, task model.Task) (model.Task, error)
	GetTaskCount(ctx context.Context) (int64, error)
	GetRunningAgentCount(ctx context.Context) (int64, error)
}
