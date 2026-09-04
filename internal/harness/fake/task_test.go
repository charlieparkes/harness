package fake

import (
	"testing"
	"time"

	"github.com/charlieparkes/go-testsize"
	"github.com/charlieparkes/harness/internal/harness/model"
)

func TestTasksReturnsRequestedCount(t *testing.T) {
	t.Parallel()
	testsize.Small(t)

	got := Tasks(20)
	if len(got) != 20 {
		t.Fatalf("Tasks(20) len = %d, want 20", len(got))
	}
	for i, task := range got {
		if task.Title == "" {
			t.Fatalf("Tasks()[%d].Title is empty", i)
		}
		if task.ID == "" {
			t.Fatalf("Tasks()[%d].ID is empty", i)
		}
		if !task.Status.IsValid() {
			t.Fatalf("Tasks()[%d].Status = %q, want a valid status", i, task.Status)
		}
	}
}

func TestTasksCyclesStatuses(t *testing.T) {
	t.Parallel()
	testsize.Small(t)

	want := []model.TaskStatus{
		model.TaskStatusReady,
		model.TaskStatusRunning,
		model.TaskStatusPendingFeedback,
		model.TaskStatusFailed,
		model.TaskStatusCompleted,
		model.TaskStatusCanceled,
	}
	got := Tasks(len(want) + 1)
	for i, task := range got {
		if task.Status != want[i%len(want)] {
			t.Fatalf("Tasks()[%d].Status = %q, want %q", i, task.Status, want[i%len(want)])
		}
	}
}

func TestTaskSetsCanceledAt(t *testing.T) {
	t.Parallel()
	testsize.Small(t)

	createdAt := time.Date(2026, 8, 30, 12, 0, 0, 0, time.UTC)
	got := Task("stop the run", createdAt, model.TaskStatusCanceled)
	if got.UpdatedAt == nil {
		t.Fatal("Task() CanceledAt = nil, want set")
	}
	want := createdAt.Add(30 * time.Minute)
	if !got.UpdatedAt.Equal(want) {
		t.Fatalf("Task() CanceledAt = %v, want %v", *got.UpdatedAt, want)
	}
}
