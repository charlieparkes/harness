package stub

import (
	"context"
	"slices"

	"github.com/charlieparkes/harness/internal/harness"
	"github.com/charlieparkes/harness/internal/harness/model"
)

func (s *Store) ListTasks(ctx context.Context) ([]model.Task, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	out := slices.Clone(s.tasks)
	slices.SortStableFunc(out, model.Task.Compare)
	return out, nil
}

func (s *Store) GetTask(ctx context.Context, id string) (model.Task, error) {
	if err := ctx.Err(); err != nil {
		return model.Task{}, err
	}

	for _, task := range s.tasks {
		if task.ID == id {
			return task, nil
		}
	}
	return model.Task{}, harness.ErrTaskNotFound
}

func (s *Store) CreateTask(ctx context.Context, task model.Task) (model.Task, error) {
	if err := ctx.Err(); err != nil {
		return model.Task{}, err
	}

	s.tasks = append(s.tasks, task)
	return task, nil
}

func (s *Store) GetTaskCount(ctx context.Context) (int64, error) {
	if err := ctx.Err(); err != nil {
		return 0, err
	}

	return int64(len(s.tasks)), nil
}
