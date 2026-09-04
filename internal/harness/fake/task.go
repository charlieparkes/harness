package fake

import (
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/charlieparkes/harness/internal/harness/model"
)

func Task(title string, createdAt time.Time, status model.TaskStatus) model.Task {
	task, err := model.NewTask(title)
	if err != nil {
		panic(err)
	}
	task.CreatedAt = createdAt
	task.Status = status
	if status == model.TaskStatusCanceled {
		canceledAt := createdAt.Add(30 * time.Minute)
		task.UpdatedAt = &canceledAt
	}
	return task
}

func Tasks(n int) []model.Task {
	faker := gofakeit.New(0)
	specs := []model.TaskStatus{
		model.TaskStatusReady,
		model.TaskStatusRunning,
		model.TaskStatusPendingFeedback,
		model.TaskStatusFailed,
		model.TaskStatusCompleted,
		model.TaskStatusCanceled,
	}
	now := time.Now().UTC()
	tasks := make([]model.Task, n)
	for i := range n {
		title := faker.Sentence()
		createdAt := faker.DateRange(now.AddDate(0, 0, -14), now).UTC()
		tasks[i] = Task(title, createdAt, specs[i%len(specs)])
	}
	return tasks
}
