package ui

import (
	"context"
	"errors"
	"slices"
	"strings"
	"testing"
	"time"

	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charlieparkes/go-testsize"
	"github.com/charlieparkes/harness/internal/harness/model"
	"github.com/charlieparkes/harness/internal/harness/ui/styles"
	"github.com/charmbracelet/colorprofile"
)

type fakeStore struct {
	tasks                  []model.Task
	err                    error
	calls                  int
	nTasks                 int64
	nRunningAgents         int64
	taskCountErr           error
	taskCountCalls         int
	runningAgentCountErr   error
	runningAgentCountCalls int
}

func (s *fakeStore) ListTasks(ctx context.Context) ([]model.Task, error) {
	s.calls++
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if s.err != nil {
		return nil, s.err
	}
	out := make([]model.Task, len(s.tasks))
	copy(out, s.tasks)
	return out, nil
}

func (s *fakeStore) GetTaskCount(ctx context.Context) (int64, error) {
	s.taskCountCalls++
	if err := ctx.Err(); err != nil {
		return 0, err
	}
	if s.taskCountErr != nil {
		return 0, s.taskCountErr
	}
	return s.nTasks, nil
}

func (s *fakeStore) GetRunningAgentCount(ctx context.Context) (int64, error) {
	s.runningAgentCountCalls++
	if err := ctx.Err(); err != nil {
		return 0, err
	}
	if s.runningAgentCountErr != nil {
		return 0, s.runningAgentCountErr
	}
	return s.nRunningAgents, nil
}

func newTasksView(ctx context.Context, store TasksStore) TasksView {
	return NewTasksView(ctx, store, styles.New(colorprofile.TrueColor))
}

func newTestTasks(ctx context.Context, t *testing.T, tasks []model.Task) TasksView {
	t.Helper()
	m := newTasksView(ctx, &fakeStore{tasks: tasks})
	return applyTasks(t, m, m.listTasksCmd()())
}

func sampleTasks() []model.Task {
	return []model.Task{
		{ID: "task-1", Title: "Draft the architecture notes"},
		{ID: "task-2", Title: "Review the open pull request"},
		{ID: "task-3", Title: "Triage failing checks"},
	}
}

func applyTasks(t *testing.T, m TasksView, msgs ...tea.Msg) TasksView {
	t.Helper()
	for _, msg := range msgs {
		m, _ = m.Update(msg)
	}
	return m
}

func selectedTask(t *testing.T, m TasksView) model.Task {
	t.Helper()

	item, ok := m.list.SelectedItem().(taskItem)
	if !ok {
		t.Fatalf("SelectedItem() = %T, want taskItem", m.list.SelectedItem())
	}
	return item.task
}

func TestTasksName(t *testing.T) {
	t.Parallel()
	testsize.Small(t)
	got := TasksView{}.Name()
	if got != "Tasks" {
		t.Fatalf("Name() = %q, want Tasks", got)
	}
}

func TestListTasksSelectsFirstTask(t *testing.T) {
	t.Parallel()
	ctx := testsize.Small(t)
	tasks := sampleTasks()
	m := newTestTasks(ctx, t, tasks)

	got := selectedTask(t, m)
	if got.ID != tasks[0].ID {
		t.Fatalf("selected ID = %q, want %q", got.ID, tasks[0].ID)
	}
}

func TestViewHidesTitle(t *testing.T) {
	t.Parallel()
	ctx := testsize.Small(t)
	m := newTestTasks(ctx, t, sampleTasks())
	if m.list.ShowTitle() {
		t.Fatal("list ShowTitle() = true, want false")
	}
	if strings.Contains(m.View(), m.Name()) {
		t.Fatalf("view %q shows list title %q", m.View(), m.Name())
	}
}

func TestViewRendersAllTasks(t *testing.T) {
	t.Parallel()
	ctx := testsize.Small(t)
	tasks := sampleTasks()
	m := newTestTasks(ctx, t, tasks)
	view := m.View()

	for _, task := range tasks {
		if !strings.Contains(view, task.Title) {
			t.Fatalf("view %q does not contain title %q", view, task.Title)
		}
	}
}

func TestTaskDescription(t *testing.T) {
	t.Parallel()
	testsize.Small(t)

	theme := styles.New(colorprofile.TrueColor)
	now := time.Date(2026, 8, 30, 12, 0, 0, 0, time.UTC)
	frame := spinner.MiniDot.Frames[0]
	tests := []struct {
		name         string
		task         model.Task
		runningFrame string
		want         string
	}{
		{
			name: "status and relative time",
			task: model.Task{
				ID:        "task-1",
				Status:    model.TaskStatusRunning,
				CreatedAt: now.Add(-2 * time.Second),
			},
			runningFrame: frame,
			want:         frame + " Running 2s ago",
		},
		{
			name: "ready day ago",
			task: model.Task{
				ID:        "task-2",
				Status:    model.TaskStatusReady,
				CreatedAt: now.Add(-24 * time.Hour),
			},
			want: "○ Ready 1d ago",
		},
		{
			name: "omits empty status and zero time",
			task: model.Task{ID: "task-3"},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := visible(taskItem{task: tt.task, theme: theme, runningFrame: tt.runningFrame}.descriptionAt(now))
			if got != tt.want {
				t.Fatalf("Description() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestTaskItemFilterValueIncludesStatus(t *testing.T) {
	t.Parallel()
	testsize.Small(t)

	task := sampleTasks()[0]
	task.Status = model.TaskStatusRunning
	got := taskItem{task: task}.FilterValue()
	want := task.Title + " " + task.ID + " " + model.TaskStatusRunning.String()
	if got != want {
		t.Fatalf("FilterValue() = %q, want %q", got, want)
	}
}

func TestViewRendersStatusAndRelativeTime(t *testing.T) {
	t.Parallel()
	ctx := testsize.Small(t)
	now := time.Now()
	task := sampleTasks()[0]
	task.Status = model.TaskStatusReady
	task.CreatedAt = now.Add(-2 * time.Second)
	tasks := []model.Task{task}
	view := newTestTasks(ctx, t, tasks).View()
	if !strings.Contains(view, model.TaskStatusReady.Text()) {
		t.Fatalf("view %q does not contain status %q", view, model.TaskStatusReady.Text())
	}
	if !strings.Contains(view, "ago") {
		t.Fatalf("view %q does not contain relative time", view)
	}
}

func TestViewStylesTitleByStatus(t *testing.T) {
	t.Parallel()
	ctx := testsize.Small(t)
	theme := styles.New(colorprofile.TrueColor)
	now := time.Now()
	tasks := []model.Task{
		{ID: "ready-1", Title: "Ready title", Status: model.TaskStatusReady, CreatedAt: now},
		{ID: "failed-1", Title: "Failed title", Status: model.TaskStatusFailed, CreatedAt: now},
	}
	m := NewTasksView(ctx, &fakeStore{tasks: tasks}, theme)
	m = applyTasks(t, m, m.listTasksCmd()())
	view := m.View()
	selected := theme.TaskItems(model.TaskStatusReady).SelectedTitle.Render(tasks[0].Title)
	if !strings.Contains(view, selected) {
		t.Fatalf("view %q does not contain ready selected title style", view)
	}
	if hasBoldSGR(selected) {
		t.Fatalf("selected title %q has bold SGR", selected)
	}
	normal := theme.TaskItems(model.TaskStatusFailed).NormalTitle.Render(tasks[1].Title)
	if !strings.Contains(view, normal) {
		t.Fatalf("view %q does not contain failed normal title style", view)
	}
	if hasBoldSGR(normal) {
		t.Fatalf("normal title %q has bold SGR", normal)
	}
}

func TestStyledTaskDescription(t *testing.T) {
	t.Parallel()
	testsize.Small(t)

	theme := styles.New(colorprofile.TrueColor)
	now := time.Date(2026, 8, 30, 12, 0, 0, 0, time.UTC)
	frame := theme.Spinner.Render(spinner.MiniDot.Frames[0])
	tests := []struct {
		name         string
		task         model.Task
		selected     bool
		runningFrame string
		want         string
	}{
		{
			name: "status and relative time",
			task: model.Task{
				ID:        "styled-1",
				Status:    model.TaskStatusRunning,
				CreatedAt: now.Add(-2 * time.Second),
			},
			runningFrame: frame,
			want:         runningStatus(theme, false, frame) + " " + mutedTimestamp(theme, "2s ago"),
		},
		{
			name: "selected running uses bright status color",
			task: model.Task{
				ID:        "styled-1-selected",
				Status:    model.TaskStatusRunning,
				CreatedAt: now.Add(-2 * time.Second),
			},
			selected:     true,
			runningFrame: frame,
			want:         runningStatus(theme, true, frame) + " " + mutedTimestamp(theme, "2s ago"),
		},
		{
			name: "ready day ago",
			task: model.Task{
				ID:        "styled-2",
				Status:    model.TaskStatusReady,
				CreatedAt: now.Add(-24 * time.Hour),
			},
			want: theme.StatusLabel(false, model.TaskStatusReady).Render(model.TaskStatusReady.Text()) + " " + mutedTimestamp(theme, relativeDayAgo),
		},
		{
			name: "omits empty status and zero time",
			task: model.Task{ID: "styled-3"},
			want: "",
		},
		{
			name: "status only",
			task: model.Task{ID: "styled-4", Status: model.TaskStatusFailed},
			want: theme.StatusLabel(false, model.TaskStatusFailed).Render(model.TaskStatusFailed.Text()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := taskItem{task: tt.task, theme: theme, selected: tt.selected, runningFrame: tt.runningFrame}.descriptionAt(now)
			if got != tt.want {
				t.Fatalf("Description() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestViewStylesDescriptionByStatus(t *testing.T) {
	t.Parallel()
	ctx := testsize.Small(t)
	theme := styles.New(colorprofile.TrueColor)
	created := time.Now().Add(-24 * time.Hour)
	tasks := []model.Task{
		{ID: "running-1", Title: "Running title", Status: model.TaskStatusRunning, CreatedAt: created},
		{ID: "ready-1", Title: "Ready title", Status: model.TaskStatusReady, CreatedAt: created},
		{ID: "feedback-1", Title: "Feedback title", Status: model.TaskStatusPendingFeedback, CreatedAt: created},
		{ID: "completed-1", Title: "Completed title", Status: model.TaskStatusCompleted, CreatedAt: created},
		{ID: "canceled-1", Title: "Canceled title", Status: model.TaskStatusCanceled, CreatedAt: created},
		{ID: "failed-1", Title: "Failed title", Status: model.TaskStatusFailed, CreatedAt: created},
		{ID: "unknown-1", Title: "Unknown title", Status: model.TaskStatus("unknown"), CreatedAt: created},
	}
	m := NewTasksView(ctx, &fakeStore{tasks: tasks}, theme)
	m = applyTasks(t, m, m.listTasksCmd()(), tea.WindowSizeMsg{Width: 80, Height: 40})
	view := m.View()
	muted := mutedTimestamp(theme, relativeDayAgo)
	if !strings.Contains(view, muted) {
		t.Fatalf("view %q does not contain muted timestamp", view)
	}
	now := time.Now()
	selectedDesc := theme.TaskItems(tasks[0].Status).SelectedDesc.Render(taskItem{task: tasks[0], theme: theme, selected: true, runningFrame: m.spinner.View()}.descriptionAt(now))
	if !strings.Contains(view, selectedDesc) {
		t.Fatalf("view %q does not contain selected description style", view)
	}
	normalDesc := theme.TaskItems(tasks[5].Status).NormalDesc.Render(taskItem{task: tasks[5], theme: theme}.descriptionAt(now))
	if !strings.Contains(view, normalDesc) {
		t.Fatalf("view %q does not contain failed normal description style", view)
	}
	for i, task := range tasks {
		if !strings.Contains(view, statusLabel(theme, i == 0, task.Status, m.spinner.View())) {
			t.Fatalf("view %q does not contain %s status label style", view, task.Status)
		}
		if task.Status == model.TaskStatusRunning {
			stale := theme.StatusLabel(i == 0, task.Status).Render(task.Status.Text())
			if strings.Contains(view, stale) {
				t.Fatalf("view %q still contains static running status text", view)
			}
		}
		items := theme.TaskItems(task.Status)
		var title string
		if i == 0 {
			title = items.SelectedTitle.Render(task.Title)
		} else {
			title = items.NormalTitle.Render(task.Title)
		}
		if !strings.Contains(view, title) {
			t.Fatalf("view %q does not contain %s title style", view, task.Status)
		}
	}
}

func TestViewDescriptionFollowsUIThemeMsg(t *testing.T) {
	t.Parallel()
	ctx := testsize.Small(t)
	initial := styles.New(colorprofile.TrueColor)
	want := styles.New(colorprofile.ANSI)
	task := model.Task{
		ID:        "theme-1",
		Title:     "Theme task",
		Status:    model.TaskStatusRunning,
		CreatedAt: time.Now().Add(-24 * time.Hour),
	}
	m := NewTasksView(ctx, &fakeStore{tasks: []model.Task{task}}, initial)
	m = applyTasks(t, m, m.listTasksCmd()(), UIThemeMsg{Theme: want})
	view := m.View()
	assertSpinnerStyle(t, m.spinner, want)
	gotFrame := m.spinner.View()
	staleFrame := initial.Spinner.Render(spinner.MiniDot.Frames[0])
	if staleFrame == gotFrame {
		t.Fatal("true-color and ANSI spinner frames match, want different palettes")
	}
	if !strings.Contains(view, gotFrame) {
		t.Fatalf("view %q does not contain themed spinner frame", view)
	}
	if strings.Contains(view, staleFrame) {
		t.Fatalf("view %q still contains stale true-color spinner frame", view)
	}
	got := want.StatusLabel(true, task.Status).Render(task.Status.String())
	stale := initial.StatusLabel(true, task.Status).Render(task.Status.String())
	if stale == got {
		t.Fatal("true-color and ANSI status labels match, want different palettes")
	}
	if !strings.Contains(view, got) {
		t.Fatalf("view %q does not contain themed status label", view)
	}
	if strings.Contains(view, stale) {
		t.Fatalf("view %q still contains stale true-color status label", view)
	}
	if !strings.Contains(view, mutedTimestamp(want, relativeDayAgo)) {
		t.Fatalf("view %q does not contain themed muted timestamp", view)
	}
	gotTitle := want.TaskItems(task.Status).SelectedTitle.Render(task.Title)
	staleTitle := initial.TaskItems(task.Status).SelectedTitle.Render(task.Title)
	if staleTitle == gotTitle {
		t.Fatal("true-color and ANSI selected titles match, want different palettes")
	}
	if !strings.Contains(view, gotTitle) {
		t.Fatalf("view %q does not contain themed selected title", view)
	}
	if strings.Contains(view, staleTitle) {
		t.Fatalf("view %q still contains stale true-color selected title", view)
	}
}

const (
	relativeDayAgo = "1d ago"
	selectedMarker = "❯"
	itemPad        = 2
)

func hasBoldSGR(s string) bool {
	return strings.Contains(s, "\x1b[1m") || strings.Contains(s, "\x1b[1;") || strings.Contains(s, ";1;") || strings.Contains(s, ";1m")
}

func mutedTimestamp(theme styles.Theme, rel string) string {
	return lipgloss.NewStyle().
		Foreground(theme.Colors.FGMuted).
		Inline(true).
		Render(rel)
}

func runningStatus(theme styles.Theme, selected bool, frame string) string {
	return frame + " " + theme.StatusLabel(selected, model.TaskStatusRunning).Render(model.TaskStatusRunning.String())
}

func statusLabel(theme styles.Theme, selected bool, status model.TaskStatus, runningFrame string) string {
	if status == model.TaskStatusRunning {
		return runningStatus(theme, selected, runningFrame)
	}
	return theme.StatusLabel(selected, status).Render(status.Text())
}

func TestDownMovesSelection(t *testing.T) {
	t.Parallel()
	ctx := testsize.Small(t)
	tasks := sampleTasks()
	m := applyTasks(t, newTestTasks(ctx, t, tasks), tea.KeyPressMsg{Text: "j"})

	got := selectedTask(t, m)
	if got.ID != tasks[1].ID {
		t.Fatalf("selected ID = %q, want %q", got.ID, tasks[1].ID)
	}
}

func TestUpMovesSelection(t *testing.T) {
	t.Parallel()
	ctx := testsize.Small(t)
	tasks := sampleTasks()
	m := newTestTasks(ctx, t, tasks)
	m = applyTasks(t, m, tea.KeyPressMsg{Text: "j"})
	m = applyTasks(t, m, tea.KeyPressMsg{Text: "k"})

	got := selectedTask(t, m)
	if got.ID != tasks[0].ID {
		t.Fatalf("selected ID = %q, want %q", got.ID, tasks[0].ID)
	}
}

func TestQuitReturnsQuitCommand(t *testing.T) {
	t.Parallel()
	ctx := testsize.Small(t)
	m := newTestTasks(ctx, t, sampleTasks())
	next, cmd := m.Update(tea.KeyPressMsg{Text: "q"})
	if cmd == nil {
		t.Fatal("Update() cmd = nil, want tea.Quit")
	}
	if _, ok := cmd().(tea.QuitMsg); !ok {
		t.Fatal("Update() cmd did not return tea.QuitMsg")
	}
	if next.View() != "" {
		t.Fatalf("quit view = %q, want empty", next.View())
	}
}

func TestInterruptQuits(t *testing.T) {
	t.Parallel()
	ctx := testsize.Small(t)
	m := newTestTasks(ctx, t, sampleTasks())
	_, cmd := m.Update(tea.InterruptMsg{})
	if cmd == nil {
		t.Fatal("Update() cmd = nil, want tea.Quit")
	}
	if _, ok := cmd().(tea.QuitMsg); !ok {
		t.Fatal("Update() cmd did not return tea.QuitMsg")
	}
}

func TestNewDoesNotListTasks(t *testing.T) {
	t.Parallel()
	ctx := testsize.Small(t)
	store := &fakeStore{tasks: sampleTasks()}
	_ = newTasksView(ctx, store)
	if store.calls != 0 {
		t.Fatalf("ListTasks() calls = %d, want 0", store.calls)
	}
}

func TestInitSchedulesListTasks(t *testing.T) {
	t.Parallel()
	ctx := testsize.Small(t)
	m := newTasksView(ctx, &fakeStore{tasks: sampleTasks()})
	if m.Init() == nil {
		t.Fatal("Init() cmd = nil, want ListTasks, refresh tick, and spinner tick")
	}
}

func TestTasksUsesMiniDot(t *testing.T) {
	t.Parallel()
	ctx := testsize.Small(t)
	m := newTasksView(ctx, &fakeStore{})
	if !slices.Equal(m.spinner.Spinner.Frames, spinner.MiniDot.Frames) {
		t.Fatalf("spinner frames = %q, want MiniDot %q", m.spinner.Spinner.Frames, spinner.MiniDot.Frames)
	}
	if m.spinner.Spinner.FPS != spinner.MiniDot.FPS {
		t.Fatalf("spinner FPS = %s, want %s", m.spinner.Spinner.FPS, spinner.MiniDot.FPS)
	}
}

func TestTasksTickAdvancesFrame(t *testing.T) {
	t.Parallel()
	ctx := testsize.Small(t)
	task := model.Task{ID: "tick-run", Title: "Running task", Status: model.TaskStatusRunning, CreatedAt: time.Now()}
	m := applyTasks(t, newTestTasks(ctx, t, []model.Task{task}), tea.WindowSizeMsg{Width: 80, Height: 24})
	first := m.spinner.View()
	if visible(first) != spinner.MiniDot.Frames[0] {
		t.Fatalf("initial spinner = %q, want %q", visible(first), spinner.MiniDot.Frames[0])
	}
	if !strings.Contains(m.View(), first) {
		t.Fatalf("view %q does not contain initial spinner frame", m.View())
	}

	m, cmd := m.Update(m.spinner.Tick())
	if cmd == nil {
		t.Fatal("first tick cmd = nil, want next spinner tick")
	}
	second := m.spinner.View()
	if visible(second) != spinner.MiniDot.Frames[1] {
		t.Fatalf("after first tick spinner = %q, want %q", visible(second), spinner.MiniDot.Frames[1])
	}
	view := m.View()
	if !strings.Contains(view, second) {
		t.Fatalf("view %q does not contain advanced spinner frame", view)
	}
	if strings.Contains(view, first) {
		t.Fatalf("view %q still contains stale spinner frame", view)
	}

	m, cmd = m.Update(m.spinner.Tick())
	if cmd == nil {
		t.Fatal("second tick cmd = nil, want next spinner tick")
	}
	if visible(m.spinner.View()) != spinner.MiniDot.Frames[2] {
		t.Fatalf("after second tick spinner = %q, want %q", visible(m.spinner.View()), spinner.MiniDot.Frames[2])
	}
}

func TestViewRunningStatusUsesThemedSpinner(t *testing.T) {
	t.Parallel()
	ctx := testsize.Small(t)
	theme := styles.New(colorprofile.TrueColor)
	now := time.Now()
	tasks := []model.Task{
		{ID: "selected-run", Title: "Selected running", Status: model.TaskStatusRunning, CreatedAt: now},
		{ID: "unselected-run", Title: "Unselected running", Status: model.TaskStatusRunning, CreatedAt: now},
		{ID: "ready-themed", Title: "Ready task", Status: model.TaskStatusReady, CreatedAt: now},
	}
	m := NewTasksView(ctx, &fakeStore{tasks: tasks}, theme)
	m = applyTasks(t, m, m.listTasksCmd()(), tea.WindowSizeMsg{Width: 80, Height: 24})
	assertSpinnerStyle(t, m.spinner, theme)
	view := m.View()
	frame := m.spinner.View()
	if !strings.Contains(view, frame) {
		t.Fatalf("view %q does not contain themed spinner frame", view)
	}
	selected := theme.StatusLabel(true, model.TaskStatusRunning).Render(model.TaskStatusRunning.String())
	unselected := theme.StatusLabel(false, model.TaskStatusRunning).Render(model.TaskStatusRunning.String())
	if selected == unselected {
		t.Fatal("selected and unselected running labels match, want different colors")
	}
	if !strings.Contains(view, selected) {
		t.Fatalf("view %q does not contain selected running label", view)
	}
	if !strings.Contains(view, unselected) {
		t.Fatalf("view %q does not contain unselected running label", view)
	}
	stale := theme.StatusLabel(true, model.TaskStatusRunning).Render(model.TaskStatusRunning.Text())
	if strings.Contains(view, stale) {
		t.Fatalf("view %q still contains static running status text", view)
	}
	if strings.Contains(visible(view), model.TaskStatusRunning.Symbol()) {
		t.Fatalf("view %q still contains static running symbol %q", view, model.TaskStatusRunning.Symbol())
	}
	if !strings.Contains(view, theme.StatusLabel(false, model.TaskStatusReady).Render(model.TaskStatusReady.Text())) {
		t.Fatalf("view %q does not contain ready status text", view)
	}
}

func TestViewSharesSpinnerFrameAcrossRunningRows(t *testing.T) {
	t.Parallel()
	ctx := testsize.Small(t)
	now := time.Now()
	tasks := []model.Task{
		{ID: "shared-run-a", Title: "First running", Status: model.TaskStatusRunning, CreatedAt: now},
		{ID: "shared-run-b", Title: "Second running", Status: model.TaskStatusRunning, CreatedAt: now},
		{ID: "shared-ready", Title: "Ready task", Status: model.TaskStatusReady, CreatedAt: now},
	}
	m := applyTasks(t, newTestTasks(ctx, t, tasks), tea.WindowSizeMsg{Width: 80, Height: 24})
	_, firstDesc := itemLines(t, m.View(), tasks[0].Title)
	_, secondDesc := itemLines(t, m.View(), tasks[1].Title)
	_, readyDesc := itemLines(t, m.View(), tasks[2].Title)
	if !strings.Contains(firstDesc, spinner.MiniDot.Frames[0]) {
		t.Fatalf("first running description %q does not contain frame %q", firstDesc, spinner.MiniDot.Frames[0])
	}
	if !strings.Contains(secondDesc, spinner.MiniDot.Frames[0]) {
		t.Fatalf("second running description %q does not contain frame %q", secondDesc, spinner.MiniDot.Frames[0])
	}
	if strings.Contains(readyDesc, spinner.MiniDot.Frames[0]) {
		t.Fatalf("ready description %q contains running spinner frame", readyDesc)
	}
	if !strings.Contains(readyDesc, model.TaskStatusReady.Text()) {
		t.Fatalf("ready description %q does not contain %q", readyDesc, model.TaskStatusReady.Text())
	}

	var cmd tea.Cmd
	m, cmd = m.Update(m.spinner.Tick())
	if cmd == nil {
		t.Fatal("tick cmd = nil, want next spinner tick")
	}
	_, firstDesc = itemLines(t, m.View(), tasks[0].Title)
	_, secondDesc = itemLines(t, m.View(), tasks[1].Title)
	if !strings.Contains(firstDesc, spinner.MiniDot.Frames[1]) {
		t.Fatalf("first running description %q does not contain frame %q", firstDesc, spinner.MiniDot.Frames[1])
	}
	if !strings.Contains(secondDesc, spinner.MiniDot.Frames[1]) {
		t.Fatalf("second running description %q does not contain frame %q", secondDesc, spinner.MiniDot.Frames[1])
	}
	if strings.Contains(firstDesc, spinner.MiniDot.Frames[0]) || strings.Contains(secondDesc, spinner.MiniDot.Frames[0]) {
		t.Fatalf("running descriptions still contain stale frame %q", spinner.MiniDot.Frames[0])
	}
}

func TestTickSchedulesListTasks(t *testing.T) {
	t.Parallel()
	ctx := testsize.Small(t)
	m := newTestTasks(ctx, t, sampleTasks())
	_, cmd := m.Update(tasksTickMsg{})
	if cmd == nil {
		t.Fatal("tick Update() cmd = nil, want ListTasks and next tick")
	}
}

func TestQuitDoesNotListTasks(t *testing.T) {
	t.Parallel()
	ctx := testsize.Small(t)
	store := &fakeStore{tasks: sampleTasks()}
	m := newTasksView(ctx, store)
	m = applyTasks(t, m, m.listTasksCmd()())
	m = applyTasks(t, m, tea.KeyPressMsg{Text: "q"})
	_, cmd := m.Update(tasksTickMsg{})
	if cmd != nil {
		t.Fatal("tick Update() after quit cmd != nil, want no ListTasks")
	}
	if store.calls != 1 {
		t.Fatalf("ListTasks() calls = %d, want 1", store.calls)
	}
}

func TestRefreshReplacesTasks(t *testing.T) {
	t.Parallel()
	ctx := testsize.Small(t)
	store := &fakeStore{tasks: sampleTasks()}
	m := newTasksView(ctx, store)
	m = applyTasks(t, m, m.listTasksCmd()())

	store.tasks = []model.Task{
		{ID: "task-9", Title: "Write the release notes"},
	}
	m = applyTasks(t, m, m.listTasksCmd()())
	view := m.View()

	if store.calls != 2 {
		t.Fatalf("ListTasks() calls = %d, want 2", store.calls)
	}
	got := selectedTask(t, m)
	if got.ID != store.tasks[0].ID {
		t.Fatalf("selected ID = %q, want %q", got.ID, store.tasks[0].ID)
	}
	if !strings.Contains(view, store.tasks[0].Title) {
		t.Fatalf("view %q did not show refreshed tasks", view)
	}
	if strings.Contains(view, sampleTasks()[0].Title) {
		t.Fatalf("view %q still shows stale tasks", view)
	}
}

func TestRefreshPreservesSelectedTask(t *testing.T) {
	t.Parallel()
	ctx := testsize.Small(t)
	store := &fakeStore{tasks: sampleTasks()}
	m := newTasksView(ctx, store)
	m = applyTasks(t, m, m.listTasksCmd()())
	m = applyTasks(t, m, tea.KeyPressMsg{Text: "j"})

	tasks := sampleTasks()
	store.tasks = []model.Task{tasks[2], tasks[1], tasks[0]}
	m = applyTasks(t, m, m.listTasksCmd()())

	got := selectedTask(t, m)
	if got.ID != tasks[1].ID {
		t.Fatalf("selected ID = %q, want %q", got.ID, tasks[1].ID)
	}
}

func TestRefreshKeepsTasksOnError(t *testing.T) {
	t.Parallel()
	ctx := testsize.Small(t)
	store := &fakeStore{
		tasks: sampleTasks(),
	}
	m := newTasksView(ctx, store)
	m = applyTasks(t, m, m.listTasksCmd()())
	store.err = errors.New("list failed")
	m = applyTasks(t, m, m.listTasksCmd()())

	if store.calls != 2 {
		t.Fatalf("ListTasks() calls = %d, want 2", store.calls)
	}
	got := selectedTask(t, m)
	if got.ID != sampleTasks()[0].ID {
		t.Fatalf("selected ID = %q, want first task after ListTasks error", got.ID)
	}
}

func TestWindowSizeResizesList(t *testing.T) {
	t.Parallel()
	ctx := testsize.Small(t)
	m := applyTasks(t, newTestTasks(ctx, t, sampleTasks()), tea.WindowSizeMsg{Width: 100, Height: 40})
	if m.list.Width() != 100 {
		t.Fatalf("list width = %d, want 100", m.list.Width())
	}
	if m.list.Height() != 40 {
		t.Fatalf("list height = %d, want 40", m.list.Height())
	}
}

func TestTasksUsesStyleGuideRoles(t *testing.T) {
	t.Parallel()
	ctx := testsize.Small(t)
	m := newTestTasks(ctx, t, sampleTasks())
	want := styles.New(colorprofile.TrueColor)
	assertTaskItemRoles(t, m, want)
	assertSpinnerStyle(t, m.spinner, want)
	plain := visible(m.View())
	if !strings.Contains(plain, selectedMarker) {
		t.Fatalf("view %q has no selected marker %q", m.View(), selectedMarker)
	}
	if strings.Contains(plain, lipgloss.NormalBorder().Left) {
		t.Fatalf("view %q still has a vertical rule", m.View())
	}
}

func TestViewSelectedMarker(t *testing.T) {
	t.Parallel()
	ctx := testsize.Small(t)
	now := time.Now()
	tasks := []model.Task{
		{ID: "selected-1", Title: "Selected architecture notes", Status: model.TaskStatusReady, CreatedAt: now},
		{ID: "normal-1", Title: "Unselected pull request", Status: model.TaskStatusFailed, CreatedAt: now},
	}
	m := applyTasks(t, newTestTasks(ctx, t, tasks), tea.WindowSizeMsg{Width: 80, Height: 24})
	view := m.View()
	plain := visible(view)
	if strings.Contains(plain, lipgloss.NormalBorder().Left) {
		t.Fatalf("view %q still has a vertical rule", view)
	}

	selTitle, selDesc := itemLines(t, view, tasks[0].Title)
	if !strings.HasPrefix(selTitle, selectedMarker) {
		t.Fatalf("selected title line %q does not start with %q", selTitle, selectedMarker)
	}
	if strings.Contains(selDesc, selectedMarker) {
		t.Fatalf("selected description %q has marker %q", selDesc, selectedMarker)
	}
	otherTitle, otherDesc := itemLines(t, view, tasks[1].Title)
	if strings.Contains(otherTitle, selectedMarker) {
		t.Fatalf("unselected title %q has marker %q", otherTitle, selectedMarker)
	}
	if strings.Contains(otherDesc, selectedMarker) {
		t.Fatalf("unselected description %q has marker %q", otherDesc, selectedMarker)
	}

	titleBefore, _, titleFound := strings.Cut(selTitle, tasks[0].Title)
	descBefore, _, descFound := strings.Cut(selDesc, tasks[0].Status.Text())
	if !titleFound || !descFound {
		t.Fatalf("selected row missing text: title %q desc %q", selTitle, selDesc)
	}
	if lipgloss.Width(titleBefore) != lipgloss.Width(descBefore) {
		t.Fatalf("selected description not aligned: title prefix %q desc prefix %q", titleBefore, descBefore)
	}
}

func TestTasksAppliesUIThemeMsg(t *testing.T) {
	t.Parallel()
	ctx := testsize.Small(t)
	want := styles.New(colorprofile.ANSI)
	m := applyTasks(t, newTestTasks(ctx, t, sampleTasks()), UIThemeMsg{Theme: want})
	assertTaskItemRoles(t, m, want)
	assertSpinnerStyle(t, m.spinner, want)
}

func TestTasksSelectedRowWidth(t *testing.T) {
	t.Parallel()
	ctx := testsize.Small(t)
	tests := []struct {
		name  string
		width int
	}{
		{name: "narrow", width: 10},
		{name: "normal", width: 80},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			m := applyTasks(t, newTestTasks(ctx, t, sampleTasks()), tea.WindowSizeMsg{Width: tc.width, Height: 24})
			for i, line := range strings.Split(m.View(), "\n") {
				if w := lipgloss.Width(line); w > tc.width {
					t.Fatalf("line %d width = %d, want <= %d", i, w, tc.width)
				}
			}
		})
	}
}

func TestUpdateAppliesAllMessages(t *testing.T) {
	t.Parallel()
	ctx := testsize.Small(t)
	tasks := sampleTasks()
	want := styles.New(colorprofile.ANSI)
	m := newTestTasks(ctx, t, tasks)
	m = applyTasks(t, m,
		tea.KeyPressMsg{Text: "j"},
		tea.KeyPressMsg{Text: "j"},
		tea.WindowSizeMsg{Width: 100, Height: 40},
		UIThemeMsg{Theme: want},
	)
	got := selectedTask(t, m)
	if got.ID != tasks[2].ID {
		t.Fatalf("selected ID = %q, want %q", got.ID, tasks[2].ID)
	}
	if m.list.Width() != 100 {
		t.Fatalf("list width = %d, want 100", m.list.Width())
	}
	if m.list.Height() != 40 {
		t.Fatalf("list height = %d, want 40", m.list.Height())
	}
	assertTaskItemRoles(t, m, want)
	assertSpinnerStyle(t, m.spinner, want)
}

func assertTaskItemRoles(t *testing.T, m TasksView, want styles.Theme) {
	t.Helper()
	p := want.Colors
	items := m.delegate.Styles
	if items.NormalTitle.GetForeground() != p.FGMuted {
		t.Fatalf("normal title foreground = %v, want fg-muted %v", items.NormalTitle.GetForeground(), p.FGMuted)
	}
	if items.NormalDesc.GetForeground() != p.FGMuted {
		t.Fatalf("normal description foreground = %v, want fg-muted %v", items.NormalDesc.GetForeground(), p.FGMuted)
	}
	if items.DimmedTitle.GetForeground() != p.FGMuted {
		t.Fatalf("dimmed title foreground = %v, want fg-muted %v", items.DimmedTitle.GetForeground(), p.FGMuted)
	}
	if items.DimmedDesc.GetForeground() != p.FGMuted {
		t.Fatalf("dimmed description foreground = %v, want fg-muted %v", items.DimmedDesc.GetForeground(), p.FGMuted)
	}
	if items.FilterMatch.GetForeground() != p.Accent {
		t.Fatalf("filter match foreground = %v, want accent %v", items.FilterMatch.GetForeground(), p.Accent)
	}
	if items.SelectedTitle.GetForeground() != p.FG {
		t.Fatalf("selected title foreground = %v, want fg %v", items.SelectedTitle.GetForeground(), p.FG)
	}
	if items.SelectedTitle.GetBold() {
		t.Fatal("selected title bold = true, want false")
	}
	if items.SelectedDesc.GetForeground() != p.FGMuted {
		t.Fatalf("selected description foreground = %v, want fg-muted %v", items.SelectedDesc.GetForeground(), p.FGMuted)
	}
	if _, ok := items.SelectedTitle.GetBackground().(lipgloss.NoColor); !ok {
		t.Fatalf("selected title background = %v, want no color", items.SelectedTitle.GetBackground())
	}
	if _, ok := items.SelectedDesc.GetBackground().(lipgloss.NoColor); !ok {
		t.Fatalf("selected description background = %v, want no color", items.SelectedDesc.GetBackground())
	}
	if !items.SelectedTitle.GetBorderLeft() || items.SelectedTitle.GetBorderLeftSize() == 0 {
		t.Fatal("selected title has no left border")
	}
	if items.SelectedTitle.GetBorderStyle().Left != selectedMarker {
		t.Fatalf("selected title glyph = %q, want %q", items.SelectedTitle.GetBorderStyle().Left, selectedMarker)
	}
	if items.SelectedTitle.GetBorderLeftForeground() != p.FG {
		t.Fatalf("selected title left border = %v, want fg %v", items.SelectedTitle.GetBorderLeftForeground(), p.FG)
	}
	if items.SelectedTitle.GetBorderLeftSize()+items.SelectedTitle.GetPaddingLeft() != itemPad {
		t.Fatalf("selected title left inset = %d, want %d", items.SelectedTitle.GetBorderLeftSize()+items.SelectedTitle.GetPaddingLeft(), itemPad)
	}
	if items.SelectedDesc.GetBorderLeft() || items.SelectedDesc.GetBorderLeftSize() != 0 {
		t.Fatal("selected description has a left border, want padding only")
	}
	if items.SelectedDesc.GetPaddingLeft() != itemPad {
		t.Fatalf("selected description padding left = %d, want %d", items.SelectedDesc.GetPaddingLeft(), itemPad)
	}
}

func itemLines(t *testing.T, view, title string) (string, string) {
	t.Helper()
	lines := strings.Split(visible(view), "\n")
	for i, line := range lines {
		if strings.Contains(line, title) {
			desc := ""
			if i+1 < len(lines) {
				desc = lines[i+1]
			}
			return line, desc
		}
	}
	t.Fatalf("view %q does not contain title %q", view, title)
	return "", ""
}
