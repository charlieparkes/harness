package storetest

import (
	"context"
	"errors"
	"reflect"
	"slices"
	"testing"
	"time"

	"github.com/charlieparkes/go-testsize"
	"github.com/charlieparkes/harness/internal/harness"
	"github.com/charlieparkes/harness/internal/harness/model"
)

// TaskStoreTests checks Task methods on a harness.Store.
// newStore must return a clean store; it is called once per subtest.
func TaskStoreTests(t *testing.T, newStore func(*testing.T) harness.Store) {
	t.Helper()
	t.Run("ListTasks", func(t *testing.T) {
		t.Parallel()

		t.Run("clean store is empty", func(t *testing.T) {
			t.Parallel()
			ctx := testsize.Small(t)
			store := newStore(t)

			got, err := store.ListTasks(ctx)
			if err != nil {
				t.Fatalf("ListTasks() error = %v", err)
			}
			if len(got) != 0 {
				t.Fatalf("ListTasks() = %#v, want empty", got)
			}
		})

		t.Run("orders by Compare", func(t *testing.T) {
			t.Parallel()
			ctx := testsize.Small(t)
			store := newStore(t)
			tasks := listOrderTasks()
			createTasks(ctx, t, store, tasks)

			got, err := store.ListTasks(ctx)
			if err != nil {
				t.Fatalf("ListTasks() error = %v", err)
			}

			want := slices.Clone(tasks)
			slices.SortStableFunc(want, model.Task.Compare)
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("ListTasks() = %#v, want %#v", got, want)
			}
		})

		t.Run("isolates returned values", func(t *testing.T) {
			t.Parallel()
			ctx := testsize.Small(t)
			store := newStore(t)
			firstTask := model.Task{ID: "iso-1", Title: "alpha", Status: model.TaskStatusReady}
			secondTask := model.Task{ID: "iso-2", Title: "beta", Status: model.TaskStatusReady}
			createTasks(ctx, t, store, []model.Task{firstTask, secondTask})

			listed, err := store.ListTasks(ctx)
			if err != nil {
				t.Fatalf("ListTasks() error = %v", err)
			}
			if len(listed) != 2 {
				t.Fatalf("ListTasks() len = %d, want 2", len(listed))
			}
			orig := listed[0]
			listed[0].Title = "mutated"

			listedAgain, err := store.ListTasks(ctx)
			if err != nil {
				t.Fatalf("second ListTasks() error = %v", err)
			}
			if !reflect.DeepEqual(listedAgain[0], orig) {
				t.Fatalf("ListTasks() returned mutated data: %#v", listedAgain[0])
			}
			if !reflect.DeepEqual(listedAgain, []model.Task{firstTask, secondTask}) {
				t.Fatalf("ListTasks() = %#v, want %#v", listedAgain, []model.Task{firstTask, secondTask})
			}

			got, err := store.GetTask(ctx, orig.ID)
			if err != nil {
				t.Fatalf("GetTask() error = %v", err)
			}
			if !reflect.DeepEqual(got, orig) {
				t.Fatalf("GetTask() after ListTasks mutation = %#v, want %#v", got, orig)
			}
		})

		t.Run("canceled context", func(t *testing.T) {
			t.Parallel()
			ctx := testsize.Small(t)
			store := newStore(t)
			createTasks(ctx, t, store, []model.Task{{ID: "listed-1", Title: "listed stored", Status: model.TaskStatusReady}})
			canceled := canceledContext(t)

			got, err := store.ListTasks(canceled)
			requireErrIs(t, err, context.Canceled)
			if got != nil {
				t.Fatalf("ListTasks() = %#v, want nil", got)
			}
		})
	})

	t.Run("GetTask", func(t *testing.T) {
		t.Parallel()

		t.Run("existing ID", func(t *testing.T) {
			t.Parallel()
			ctx := testsize.Small(t)
			store := newStore(t)
			want := model.Task{ID: "get-2", Title: "looked up", Status: model.TaskStatusRunning, CreatedAt: time.Date(2026, 8, 27, 11, 0, 0, 0, time.UTC)}
			createTasks(ctx, t, store, []model.Task{
				{ID: "get-1", Title: "other", Status: model.TaskStatusReady},
				want,
			})

			got, err := store.GetTask(ctx, want.ID)
			if err != nil {
				t.Fatalf("GetTask() error = %v", err)
			}
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("GetTask() = %#v, want %#v", got, want)
			}
		})

		t.Run("missing ID", func(t *testing.T) {
			t.Parallel()
			ctx := testsize.Small(t)
			store := newStore(t)
			createTasks(ctx, t, store, []model.Task{{ID: "exists", Title: "here", Status: model.TaskStatusReady}})

			got, err := store.GetTask(ctx, "missing")
			requireErrIs(t, err, harness.ErrTaskNotFound)
			if !reflect.DeepEqual(got, model.Task{}) {
				t.Fatalf("GetTask() = %#v, want zero Task", got)
			}
		})

		t.Run("isolates returned values", func(t *testing.T) {
			t.Parallel()
			ctx := testsize.Small(t)
			store := newStore(t)
			want := model.Task{ID: "get-copy-1", Title: "original", Status: model.TaskStatusReady}
			createTasks(ctx, t, store, []model.Task{want})

			got, err := store.GetTask(ctx, want.ID)
			if err != nil {
				t.Fatalf("GetTask() error = %v", err)
			}
			got.Title = "mutated"

			second, err := store.GetTask(ctx, want.ID)
			if err != nil {
				t.Fatalf("second GetTask() error = %v", err)
			}
			if !reflect.DeepEqual(second, want) {
				t.Fatalf("GetTask() returned mutated data: %#v", second)
			}
		})

		t.Run("canceled context", func(t *testing.T) {
			t.Parallel()
			ctx := testsize.Small(t)
			store := newStore(t)
			task := model.Task{ID: "get-cancel-1", Title: "get stored", Status: model.TaskStatusReady}
			createTasks(ctx, t, store, []model.Task{task})
			canceled := canceledContext(t)

			got, err := store.GetTask(canceled, task.ID)
			requireErrIs(t, err, context.Canceled)
			if !reflect.DeepEqual(got, model.Task{}) {
				t.Fatalf("GetTask() = %#v, want zero Task", got)
			}
		})
	})

	t.Run("CreateTask", func(t *testing.T) {
		t.Parallel()

		t.Run("stores and returns the task", func(t *testing.T) {
			t.Parallel()
			ctx := testsize.Small(t)
			store := newStore(t)
			existing := model.Task{ID: "existing-1", Title: "existing", Status: model.TaskStatusReady}
			want := model.Task{ID: "created-1", Title: "created", Status: model.TaskStatusReady, CreatedAt: time.Date(2026, 8, 27, 12, 0, 0, 0, time.UTC)}
			createTasks(ctx, t, store, []model.Task{existing})

			got, err := store.CreateTask(ctx, want)
			if err != nil {
				t.Fatalf("CreateTask() error = %v", err)
			}
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("CreateTask() = %#v, want %#v", got, want)
			}

			stored, err := store.GetTask(ctx, want.ID)
			if err != nil {
				t.Fatalf("GetTask() error = %v", err)
			}
			if !reflect.DeepEqual(stored, want) {
				t.Fatalf("GetTask() after CreateTask() = %#v, want %#v", stored, want)
			}

			listed, err := store.ListTasks(ctx)
			if err != nil {
				t.Fatalf("ListTasks() error = %v", err)
			}
			if !reflect.DeepEqual(listed, []model.Task{want, existing}) {
				t.Fatalf("ListTasks() after CreateTask() = %#v, want %#v", listed, []model.Task{want, existing})
			}
		})

		t.Run("isolates input and returned values", func(t *testing.T) {
			t.Parallel()
			ctx := testsize.Small(t)
			store := newStore(t)
			task := model.Task{ID: "copy-1", Title: "kept", Status: model.TaskStatusReady}

			got, err := store.CreateTask(ctx, task)
			if err != nil {
				t.Fatalf("CreateTask() error = %v", err)
			}
			got.Title = "changed-return"
			task.Title = "changed-input"

			stored, err := store.GetTask(ctx, "copy-1")
			if err != nil {
				t.Fatalf("GetTask() error = %v", err)
			}
			if stored.Title != "kept" {
				t.Fatalf("CreateTask() stored mutated data: %#v", stored)
			}
		})

		t.Run("canceled context", func(t *testing.T) {
			t.Parallel()
			ctx := testsize.Small(t)
			store := newStore(t)
			canceled := canceledContext(t)

			got, err := store.CreateTask(canceled, model.Task{ID: "canceled-1", Title: "new", Status: model.TaskStatusReady})
			requireErrIs(t, err, context.Canceled)
			if !reflect.DeepEqual(got, model.Task{}) {
				t.Fatalf("CreateTask() = %#v, want zero Task", got)
			}

			listed, err := store.ListTasks(ctx)
			if err != nil {
				t.Fatalf("ListTasks() error = %v", err)
			}
			if len(listed) != 0 {
				t.Fatalf("CreateTask() stored %#v, want empty", listed)
			}
		})
	})

	t.Run("GetTaskCount", func(t *testing.T) {
		t.Parallel()

		t.Run("clean store is zero", func(t *testing.T) {
			t.Parallel()
			ctx := testsize.Small(t)
			store := newStore(t)

			got, err := store.GetTaskCount(ctx)
			if err != nil {
				t.Fatalf("GetTaskCount() error = %v", err)
			}
			if got != 0 {
				t.Fatalf("GetTaskCount() = %d, want 0", got)
			}
		})

		t.Run("increases after each create", func(t *testing.T) {
			t.Parallel()
			ctx := testsize.Small(t)
			store := newStore(t)
			tasks := []model.Task{
				{ID: "count-1", Title: "one", Status: model.TaskStatusReady},
				{ID: "count-2", Title: "two", Status: model.TaskStatusReady},
				{ID: "count-3", Title: "three", Status: model.TaskStatusReady},
			}

			for i, task := range tasks {
				if _, err := store.CreateTask(ctx, task); err != nil {
					t.Fatalf("CreateTask() error = %v", err)
				}
				got, err := store.GetTaskCount(ctx)
				if err != nil {
					t.Fatalf("GetTaskCount() error = %v", err)
				}
				want := int64(i + 1)
				if got != want {
					t.Fatalf("GetTaskCount() = %d, want %d", got, want)
				}
			}
		})

		t.Run("canceled context", func(t *testing.T) {
			t.Parallel()
			ctx := testsize.Small(t)
			store := newStore(t)
			createTasks(ctx, t, store, []model.Task{{ID: "count-1", Title: "count stored", Status: model.TaskStatusReady}})
			canceled := canceledContext(t)

			got, err := store.GetTaskCount(canceled)
			requireErrIs(t, err, context.Canceled)
			if got != 0 {
				t.Fatalf("GetTaskCount() = %d, want 0", got)
			}
		})
	})
}

func listOrderTasks() []model.Task {
	newer := time.Date(2026, 8, 27, 12, 0, 0, 0, time.UTC)
	older := time.Date(2026, 8, 27, 10, 0, 0, 0, time.UTC)
	canceledAt := newer
	return []model.Task{
		{ID: "completed-newer", Title: "completed newer", Status: model.TaskStatusCompleted, CreatedAt: newer},
		{ID: "ready-older", Title: "ready older", Status: model.TaskStatusReady, CreatedAt: older},
		{ID: "in-progress-newer", Title: "in progress newer", Status: model.TaskStatusRunning, CreatedAt: newer},
		{ID: "feedback-older", Title: "feedback older", Status: model.TaskStatusPendingFeedback, CreatedAt: older},
		{ID: "failed-newer", Title: "failed newer", Status: model.TaskStatusFailed, CreatedAt: newer},
		{ID: "canceled-status", Title: "canceled status", Status: model.TaskStatusCanceled, CreatedAt: newer},
		{ID: "ready-canceled", Title: "ready canceled", Status: model.TaskStatusReady, CreatedAt: newer, UpdatedAt: &canceledAt},
		{ID: "completed-older", Title: "completed older", Status: model.TaskStatusCompleted, CreatedAt: older},
		{ID: "in-progress-tie", Title: "in progress tie", Status: model.TaskStatusRunning, CreatedAt: newer},
		{ID: "completed-tie", Title: "completed tie", Status: model.TaskStatusCompleted, CreatedAt: newer},
	}
}

func createTasks(ctx context.Context, t *testing.T, store harness.Store, tasks []model.Task) {
	t.Helper()
	for _, task := range tasks {
		if _, err := store.CreateTask(ctx, task); err != nil {
			t.Fatalf("CreateTask() error = %v", err)
		}
	}
}

func canceledContext(t *testing.T) context.Context {
	t.Helper()
	ctx, cancel := context.WithCancel(testsize.Small(t))
	cancel()
	return ctx
}

func requireErrIs(t *testing.T, err, want error) {
	t.Helper()
	if err == nil {
		t.Fatalf("error = nil, want %v", want)
	}
	if !errors.Is(err, want) {
		t.Fatalf("error = %v, want %v", err, want)
	}
}
