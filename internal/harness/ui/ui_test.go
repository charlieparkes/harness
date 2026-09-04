package ui

import (
	"context"
	"regexp"
	"strings"
	"testing"

	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charlieparkes/go-testsize"
	"github.com/charlieparkes/harness/internal/harness/model"
	"github.com/charlieparkes/harness/internal/harness/ui/styles"
	"github.com/charmbracelet/colorprofile"
)

const (
	keyTab      = "tab"
	keyShiftTab = "shift+tab"
	keyLeft     = "left"
	keyRight    = "right"
	keyDown     = "j"
	keyQuit     = "q"
)

var ansiEscape = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func tasksOf(t *testing.T, m Model) TasksView {
	t.Helper()
	return m.tasks
}

func exampleOf(t *testing.T, m Model) ExampleView {
	t.Helper()
	return m.example
}

func statusOf(t *testing.T, m Model) StatusView {
	t.Helper()
	return m.status
}

func apply(t *testing.T, m Model, msg tea.Msg) Model {
	t.Helper()

	next, _ := m.Update(msg)
	model, ok := next.(Model)
	if !ok {
		t.Fatalf("Update() returned %T, want Model", next)
	}
	return model
}

func newTestModel(ctx context.Context, t *testing.T, tasks []model.Task) Model {
	t.Helper()
	m := New(ctx, &fakeStore{tasks: tasks})
	return apply(t, m, tasksOf(t, m).listTasksCmd()())
}

func visible(s string) string {
	return ansiEscape.ReplaceAllString(s, "")
}

func statusLine(t *testing.T, m Model) string {
	t.Helper()
	view := m.render()
	line, _, ok := strings.Cut(view, "\n")
	if !ok {
		t.Fatalf("view %q has no status bar line", view)
	}
	return visible(line)
}

func TestDefaultViewIsTasks(t *testing.T) {
	t.Parallel()
	ctx := testsize.Small(t)
	m := newTestModel(ctx, t, sampleTasks())
	if m.activeView != ViewTypeTasks {
		t.Fatalf("activeView = %d, want %d", m.activeView, ViewTypeTasks)
	}
	view := m.render()
	if !strings.Contains(view, sampleTasks()[0].Title) {
		t.Fatalf("view %q does not show tasks on the default view", view)
	}
}

func TestStatusBarShowsActiveViewName(t *testing.T) {
	t.Parallel()
	ctx := testsize.Small(t)
	m := newTestModel(ctx, t, sampleTasks())
	line := statusLine(t, m)
	want := programName + breadcrumbSep + m.tasks.Name()
	if !strings.Contains(line, want) {
		t.Fatalf("status bar %q does not contain breadcrumb %q", line, want)
	}

	m = m.setActiveView(ViewTypeExample)
	line = statusLine(t, m)
	want = programName + breadcrumbSep + m.example.Name()
	if !strings.Contains(line, want) {
		t.Fatalf("status bar %q does not contain breadcrumb %q", line, want)
	}
}

func TestStatusBarShowsStoreStats(t *testing.T) {
	t.Parallel()
	ctx := testsize.Small(t)
	store := &fakeStore{
		tasks:          sampleTasks(),
		nTasks:         9,
		nRunningAgents: 4,
	}
	m := New(ctx, store)
	m = apply(t, m, m.status.statsCmd()())
	line := statusLine(t, m)
	want := "9 tasks, 4 agents running"
	if !strings.Contains(line, want) {
		t.Fatalf("status bar %q does not contain stats %q", line, want)
	}
}

func TestRenderDoesNotFetchStats(t *testing.T) {
	t.Parallel()
	ctx := testsize.Small(t)
	store := &fakeStore{}
	m := New(ctx, store)
	_ = m.render()
	_ = m.render()
	if store.taskCountCalls != 0 {
		t.Fatalf("GetTaskCount() calls = %d, want 0", store.taskCountCalls)
	}
	if store.runningAgentCountCalls != 0 {
		t.Fatalf("GetRunningAgentCount() calls = %d, want 0", store.runningAgentCountCalls)
	}
}

func TestStatusBarDrawsBorderAboveActiveView(t *testing.T) {
	t.Parallel()
	ctx := testsize.Small(t)
	m := newTestModel(ctx, t, sampleTasks())
	rule := strings.Repeat(lipgloss.NormalBorder().Bottom, m.width)
	viewLines := strings.Split(visible(m.render()), "\n")
	if len(viewLines) < 3 {
		t.Fatalf("view lines = %d, want at least 3", len(viewLines))
	}
	if viewLines[1] != rule {
		t.Fatalf("view border = %q, want %q", viewLines[1], rule)
	}
	if !strings.Contains(strings.Join(viewLines[2:], "\n"), sampleTasks()[0].Title) {
		t.Fatalf("view %q does not show the active view below the border", m.render())
	}
}

func TestRenderShowsOnlyActiveView(t *testing.T) {
	t.Parallel()
	ctx := testsize.Small(t)
	m := newTestModel(ctx, t, sampleTasks())
	view := m.render()
	if strings.Contains(view, m.example.View()) {
		t.Fatalf("view %q shows the inactive Example view", view)
	}
	if !strings.Contains(view, sampleTasks()[0].Title) {
		t.Fatalf("view %q does not show the active Tasks view", view)
	}

	m = m.setActiveView(ViewTypeExample)
	view = m.render()
	if strings.Contains(view, sampleTasks()[0].Title) {
		t.Fatalf("view %q still shows Tasks after switching the active view", view)
	}
	if !strings.Contains(view, m.example.View()) {
		t.Fatalf("view %q does not show the Example view", view)
	}
	line := statusLine(t, m)
	want := programName + breadcrumbSep + m.example.Name()
	if !strings.HasPrefix(strings.TrimLeft(line, " "), want) {
		t.Fatalf("status bar %q does not start with %q", line, want)
	}
}

func TestShellKeysDoNotSwitchViews(t *testing.T) {
	t.Parallel()
	ctx := testsize.Small(t)
	m := newTestModel(ctx, t, sampleTasks())
	for _, key := range []string{keyTab, keyShiftTab, keyLeft, keyRight} {
		m = apply(t, m, tea.KeyPressMsg{Text: key})
		if m.activeView != ViewTypeTasks {
			t.Fatalf("after %q, activeView = %d, want %d", key, m.activeView, ViewTypeTasks)
		}
		if !strings.Contains(m.render(), sampleTasks()[0].Title) {
			t.Fatalf("after %q, view %q does not show Tasks", key, m.render())
		}
	}
}

func TestActiveViewReceivesKeys(t *testing.T) {
	t.Parallel()
	ctx := testsize.Small(t)
	tasks := sampleTasks()
	m := apply(t, newTestModel(ctx, t, tasks), tea.KeyPressMsg{Text: keyDown})
	got := selectedTask(t, tasksOf(t, m))
	if got.ID != tasks[1].ID {
		t.Fatalf("selected ID = %q, want %q", got.ID, tasks[1].ID)
	}
}

func TestHiddenViewDoesNotReceiveKeys(t *testing.T) {
	t.Parallel()
	ctx := testsize.Small(t)
	tasks := sampleTasks()
	m := newTestModel(ctx, t, tasks)
	m = m.setActiveView(ViewTypeExample)
	m = apply(t, m, tea.KeyPressMsg{Text: keyDown})
	got := selectedTask(t, tasksOf(t, m))
	if got.ID != tasks[0].ID {
		t.Fatalf("hidden Tasks selected ID = %q, want %q", got.ID, tasks[0].ID)
	}
}

func TestWindowSizeResizesActiveChild(t *testing.T) {
	t.Parallel()
	ctx := testsize.Small(t)
	m := apply(t, newTestModel(ctx, t, sampleTasks()), tea.WindowSizeMsg{Width: 100, Height: 40})
	wantWidth, wantHeight := m.childWindowSize()
	if m.width != 100 {
		t.Fatalf("width = %d, want 100", m.width)
	}
	if m.height != 40 {
		t.Fatalf("height = %d, want 40", m.height)
	}
	tasks := tasksOf(t, m)
	if tasks.list.Width() != wantWidth {
		t.Fatalf("tasks width = %d, want %d", tasks.list.Width(), wantWidth)
	}
	if tasks.list.Height() != wantHeight {
		t.Fatalf("tasks height = %d, want %d", tasks.list.Height(), wantHeight)
	}
	if statusOf(t, m).width != 100 {
		t.Fatalf("status width = %d, want 100", statusOf(t, m).width)
	}
	if wantHeight >= 40 {
		t.Fatalf("child height = %d, want less than window height 40", wantHeight)
	}

	m = m.setActiveView(ViewTypeExample)
	m = apply(t, m, tea.WindowSizeMsg{Width: 80, Height: 30})
	wantWidth, wantHeight = m.childWindowSize()
	example := exampleOf(t, m)
	if example.width != wantWidth {
		t.Fatalf("example width = %d, want %d", example.width, wantWidth)
	}
	if example.height != wantHeight {
		t.Fatalf("example height = %d, want %d", example.height, wantHeight)
	}
	if statusOf(t, m).width != 80 {
		t.Fatalf("status width = %d, want 80", statusOf(t, m).width)
	}
}

func TestWindowSizeResizesHiddenViews(t *testing.T) {
	t.Parallel()
	ctx := testsize.Small(t)
	m := apply(t, newTestModel(ctx, t, sampleTasks()), tea.WindowSizeMsg{Width: 100, Height: 40})
	wantWidth, wantHeight := m.childWindowSize()
	example := exampleOf(t, m)
	if example.width != wantWidth {
		t.Fatalf("hidden example width = %d, want %d", example.width, wantWidth)
	}
	if example.height != wantHeight {
		t.Fatalf("hidden example height = %d, want %d", example.height, wantHeight)
	}

	m = m.setActiveView(ViewTypeExample)
	m = apply(t, m, tea.WindowSizeMsg{Width: 80, Height: 30})
	wantWidth, wantHeight = m.childWindowSize()
	tasks := tasksOf(t, m)
	if tasks.list.Width() != wantWidth {
		t.Fatalf("hidden tasks width = %d, want %d", tasks.list.Width(), wantWidth)
	}
	if tasks.list.Height() != wantHeight {
		t.Fatalf("hidden tasks height = %d, want %d", tasks.list.Height(), wantHeight)
	}
}

func TestHiddenTasksStillRefresh(t *testing.T) {
	t.Parallel()
	ctx := testsize.Small(t)
	store := &fakeStore{tasks: sampleTasks()}
	m := New(ctx, store)
	m = apply(t, m, tasksOf(t, m).listTasksCmd()())
	m = m.setActiveView(ViewTypeExample)

	store.tasks = []model.Task{
		{ID: "task-9", Title: "Write the release notes"},
	}
	m = apply(t, m, tasksOf(t, m).listTasksCmd()())
	m = m.setActiveView(ViewTypeTasks)

	if store.calls != 2 {
		t.Fatalf("ListTasks() calls = %d, want 2", store.calls)
	}
	if !strings.Contains(m.render(), "Write the release notes") {
		t.Fatalf("view %q did not show tasks refreshed while hidden", m.render())
	}
}

func TestHiddenExampleStillSpins(t *testing.T) {
	t.Parallel()
	ctx := testsize.Small(t)
	m := newTestModel(ctx, t, sampleTasks())
	before := m.example.View()
	m = apply(t, m, exampleOf(t, m).spinner.Tick())
	if m.activeView != ViewTypeTasks {
		t.Fatalf("activeView = %d, want %d", m.activeView, ViewTypeTasks)
	}
	if m.example.View() == before {
		t.Fatal("hidden Example spinner did not advance")
	}
}

func TestInitStartsChildren(t *testing.T) {
	t.Parallel()
	ctx := testsize.Small(t)
	m := New(ctx, &fakeStore{tasks: sampleTasks()})
	if m.Init() == nil {
		t.Fatal("Init() cmd = nil, want Tasks refresh and Example spinner tick")
	}
}

func TestShellQuitReturnsQuitCommand(t *testing.T) {
	t.Parallel()
	ctx := testsize.Small(t)
	m := newTestModel(ctx, t, sampleTasks())
	next, cmd := m.Update(tea.KeyPressMsg{Text: keyQuit})
	model, ok := next.(Model)
	if !ok {
		t.Fatalf("Update() returned %T, want Model", next)
	}
	if cmd == nil {
		t.Fatal("Update() cmd = nil, want tea.Quit")
	}
	if _, ok := cmd().(tea.QuitMsg); !ok {
		t.Fatal("Update() cmd did not return tea.QuitMsg")
	}
	if model.render() != "" {
		t.Fatalf("quit view = %q, want empty", model.render())
	}
}

func TestQuitFromExampleView(t *testing.T) {
	t.Parallel()
	ctx := testsize.Small(t)
	m := newTestModel(ctx, t, sampleTasks())
	m = m.setActiveView(ViewTypeExample)
	next, cmd := m.Update(tea.KeyPressMsg{Text: keyQuit})
	model, ok := next.(Model)
	if !ok {
		t.Fatalf("Update() returned %T, want Model", next)
	}
	if cmd == nil {
		t.Fatal("Update() cmd = nil, want tea.Quit")
	}
	if _, ok := cmd().(tea.QuitMsg); !ok {
		t.Fatal("Update() cmd did not return tea.QuitMsg")
	}
	if model.render() != "" {
		t.Fatalf("quit view = %q, want empty", model.render())
	}
}

func TestShellInterruptQuits(t *testing.T) {
	t.Parallel()
	ctx := testsize.Small(t)
	m := newTestModel(ctx, t, sampleTasks())
	_, cmd := m.Update(tea.InterruptMsg{})
	if cmd == nil {
		t.Fatal("Update() cmd = nil, want tea.Quit")
	}
	if _, ok := cmd().(tea.QuitMsg); !ok {
		t.Fatal("Update() cmd did not return tea.QuitMsg")
	}
}

func TestShellQuitDoesNotListTasks(t *testing.T) {
	t.Parallel()
	ctx := testsize.Small(t)
	store := &fakeStore{tasks: sampleTasks()}
	m := New(ctx, store)
	m = apply(t, m, tasksOf(t, m).listTasksCmd()())
	m = apply(t, m, tea.KeyPressMsg{Text: keyQuit})
	_, cmd := m.Update(tasksTickMsg{})
	if cmd != nil {
		t.Fatal("tick Update() after quit cmd != nil, want no ListTasks")
	}
	if store.calls != 1 {
		t.Fatalf("ListTasks() calls = %d, want 1", store.calls)
	}
}

func TestNewDoesNotListTasksFromShell(t *testing.T) {
	t.Parallel()
	ctx := testsize.Small(t)
	store := &fakeStore{tasks: sampleTasks()}
	_ = New(ctx, store)
	if store.calls != 0 {
		t.Fatalf("ListTasks() calls = %d, want 0", store.calls)
	}
}

func TestMiniDotFrameAppearsOnExampleView(t *testing.T) {
	t.Parallel()
	ctx := testsize.Small(t)
	m := newTestModel(ctx, t, sampleTasks())
	m = m.setActiveView(ViewTypeExample)
	if !strings.Contains(m.render(), spinner.MiniDot.Frames[0]) {
		t.Fatalf("view %q does not contain MiniDot frame %q", m.render(), spinner.MiniDot.Frames[0])
	}
}

func TestColorProfileRebuildsTheme(t *testing.T) {
	t.Parallel()
	ctx := testsize.Small(t)
	m := New(ctx, &fakeStore{})

	for _, profile := range []colorprofile.Profile{
		colorprofile.ANSI,
		colorprofile.TrueColor,
	} {
		next, cmd := m.Update(tea.ColorProfileMsg{Profile: profile})
		if cmd != nil {
			t.Fatalf("ColorProfileMsg %s cmd != nil, want nil", profile)
		}
		var ok bool
		m, ok = next.(Model)
		if !ok {
			t.Fatalf("Update() returned %T, want Model", next)
		}
		want := styles.New(profile)
		assertThemeApplied(t, m, want)
		assertStatusRoles(t, statusOf(t, m), want)
		assertTaskItemRoles(t, tasksOf(t, m), want)
		assertSpinnerStyle(t, exampleOf(t, m).spinner, want)
		assertSpinnerStyle(t, tasksOf(t, m).spinner, want)
	}
}

func TestColorProfileUpdatesAllViews(t *testing.T) {
	t.Parallel()
	ctx := testsize.Small(t)
	m := New(ctx, &fakeStore{})
	m = m.setActiveView(ViewTypeExample)
	want := styles.New(colorprofile.ANSI)
	m = apply(t, m, tea.ColorProfileMsg{Profile: colorprofile.ANSI})
	assertSpinnerStyle(t, exampleOf(t, m).spinner, want)
	assertSpinnerStyle(t, tasksOf(t, m).spinner, want)
	assertStatusRoles(t, statusOf(t, m), want)
	assertTaskItemRoles(t, tasksOf(t, m), want)
}

func assertThemeApplied(t *testing.T, m Model, want styles.Theme) {
	t.Helper()
	if m.theme.Colors != want.Colors {
		t.Fatalf("theme colors = %+v, want %+v", m.theme.Colors, want.Colors)
	}
	view := m.View()
	if view.ForegroundColor != nil {
		t.Fatalf("canvas foreground = %v, want unset", view.ForegroundColor)
	}
	if view.BackgroundColor != nil {
		t.Fatalf("canvas background = %v, want unset", view.BackgroundColor)
	}
}
