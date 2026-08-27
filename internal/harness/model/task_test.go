package model

import (
	"strings"
	"testing"
	"time"

	"github.com/charlieparkes/go-testsize"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

func TestNewTaskSetsTitleAndNanoid(t *testing.T) {
	t.Parallel()
	testsize.Small(t)

	before := time.Now().UTC()
	task, err := NewTask("Draft the architecture notes")
	after := time.Now().UTC()
	if err != nil {
		t.Fatalf("NewTask() error = %v", err)
	}

	if task.Title != "Draft the architecture notes" {
		t.Fatalf("NewTask() Title = %q, want %q", task.Title, "Draft the architecture notes")
	}
	if len(task.ID) != 4 {
		t.Fatalf("NewTask() ID length = %d, want 4", len(task.ID))
	}
	for _, r := range task.ID {
		if !strings.ContainsRune(gonanoid.CrockfordBase32Upper, r) {
			t.Fatalf("NewTask() ID %q contains %q, want Crockford Base32", task.ID, r)
		}
	}
	if task.Status != TaskStatusReady {
		t.Fatalf("NewTask() Status = %q, want %q", task.Status, TaskStatusReady)
	}
	if task.UpdatedAt != nil {
		t.Fatalf("NewTask() CanceledAt = %v, want nil", task.UpdatedAt)
	}
	if task.CreatedAt.Location() != time.UTC {
		t.Fatalf("NewTask() CreatedAt location = %s, want UTC", task.CreatedAt.Location())
	}
	if task.CreatedAt.Before(before) || task.CreatedAt.After(after) {
		t.Fatalf("NewTask() CreatedAt = %v, want between %v and %v", task.CreatedAt, before, after)
	}
}

func TestNewTaskIDsAreUnique(t *testing.T) {
	t.Parallel()
	testsize.Small(t)

	first, err := NewTask("one")
	if err != nil {
		t.Fatalf("first NewTask() error = %v", err)
	}

	second, err := NewTask("two")
	if err != nil {
		t.Fatalf("second NewTask() error = %v", err)
	}

	if first.ID == second.ID {
		t.Fatalf("NewTask() generated duplicate ID %q", first.ID)
	}
}

func TestTaskIsActive(t *testing.T) {
	t.Parallel()
	testsize.Small(t)

	tests := []struct {
		name string
		task Task
		want bool
	}{
		{name: "ready", task: Task{Status: TaskStatusReady}, want: true},
		{name: "running", task: Task{Status: TaskStatusRunning}, want: true},
		{name: "pending_feedback", task: Task{Status: TaskStatusPendingFeedback}, want: true},
		{name: "failed", task: Task{Status: TaskStatusFailed}, want: true},
		{name: "completed", task: Task{Status: TaskStatusCompleted}, want: false},
		{name: "canceled", task: Task{Status: TaskStatusCanceled}, want: false},
		{name: "unknown status", task: Task{Status: TaskStatus("unknown")}, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.task.IsActive(); got != tt.want {
				t.Fatalf("IsActive() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTaskCompare(t *testing.T) {
	t.Parallel()
	testsize.Small(t)

	newer := time.Date(2026, 8, 27, 12, 0, 0, 0, time.UTC)
	older := time.Date(2026, 8, 27, 10, 0, 0, 0, time.UTC)
	latest := time.Date(2026, 8, 27, 14, 0, 0, 0, time.UTC)

	activeNewer := Task{Status: TaskStatusReady, CreatedAt: newer}
	activeOlder := Task{Status: TaskStatusRunning, CreatedAt: older}
	inactiveNewer := Task{Status: TaskStatusCompleted, CreatedAt: newer}
	inactiveOlder := Task{Status: TaskStatusCanceled, CreatedAt: older}
	activeUpdated := Task{Status: TaskStatusFailed, CreatedAt: older, UpdatedAt: &latest}
	inactiveUpdated := Task{Status: TaskStatusCompleted, CreatedAt: older, UpdatedAt: &latest}

	tests := []struct {
		name  string
		left  Task
		right Task
		want  int
	}{
		{name: "active before inactive", left: activeOlder, right: inactiveNewer, want: -1},
		{name: "inactive after active", left: inactiveNewer, right: activeOlder, want: 1},
		{name: "newer active before older active", left: activeNewer, right: activeOlder, want: -1},
		{name: "older active after newer active", left: activeOlder, right: activeNewer, want: 1},
		{name: "updated activity before created-only", left: activeUpdated, right: activeNewer, want: -1},
		{name: "newer inactive before older inactive", left: inactiveNewer, right: inactiveOlder, want: -1},
		{name: "updated inactive before created-only inactive", left: inactiveUpdated, right: inactiveNewer, want: -1},
		{name: "equal active timestamps", left: activeNewer, right: Task{Status: TaskStatusFailed, CreatedAt: newer}, want: 0},
		{name: "equal inactive timestamps", left: inactiveNewer, right: Task{Status: TaskStatusCanceled, CreatedAt: newer}, want: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.left.Compare(tt.right); got != tt.want {
				t.Fatalf("Compare() = %d, want %d", got, tt.want)
			}
		})
	}
}
