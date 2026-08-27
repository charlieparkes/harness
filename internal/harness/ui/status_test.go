package ui

import (
	"context"
	"errors"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charlieparkes/go-testsize"
	"github.com/charlieparkes/harness/internal/harness/ui/styles"
	"github.com/charmbracelet/colorprofile"
)

const fakeStats = "3 tasks, 2 agents running"

func newStatusView(ctx context.Context) StatusView {
	m := NewStatusView(ctx, &fakeStore{nTasks: 3, nRunningAgents: 2}, styles.New(colorprofile.TrueColor))
	m, _ = m.Update(StatusViewActiveViewMsg{View: ViewTypeTasks})
	m, _ = m.Update(m.statsCmd()())
	return m
}

func applyStatus(t *testing.T, m StatusView, msgs ...tea.Msg) StatusView {
	t.Helper()
	for _, msg := range msgs {
		m, _ = m.Update(msg)
	}
	return m
}

func statusViewLine(t *testing.T, m StatusView) string {
	t.Helper()
	view := m.View()
	line, _, ok := strings.Cut(view, "\n")
	if !ok {
		t.Fatalf("status view %q has no status bar line", view)
	}
	return visible(line)
}

func TestStatusShowsTitleAndStats(t *testing.T) {
	t.Parallel()
	ctx := testsize.Small(t)
	m := newStatusView(ctx)
	line := statusViewLine(t, m)
	want := programName + breadcrumbSep + tasksTitle
	if !strings.Contains(line, want) {
		t.Fatalf("status bar %q does not contain breadcrumb %q", line, want)
	}
	if !strings.Contains(line, fakeStats) {
		t.Fatalf("status bar %q does not contain fake stats %q", line, fakeStats)
	}
}

func TestStatusUsesGivenView(t *testing.T) {
	t.Parallel()
	ctx := testsize.Small(t)
	m := applyStatus(t, newStatusView(ctx), StatusViewActiveViewMsg{View: ViewTypeExample})
	if m.activeView != ViewTypeExample {
		t.Fatalf("active view = %v, want %v", m.activeView, ViewTypeExample)
	}
	want := programName + breadcrumbSep + ViewTypeExample.Name()
	line := statusViewLine(t, m)
	if !strings.HasPrefix(strings.TrimLeft(line, " "), want) {
		t.Fatalf("status bar %q does not start with %q", line, want)
	}
}

func TestStatusRootOnlyBreadcrumb(t *testing.T) {
	t.Parallel()
	ctx := testsize.Small(t)
	invalidView := ViewType(-1)
	m := applyStatus(t, newStatusView(ctx), StatusViewActiveViewMsg{View: invalidView})
	if m.activeView != invalidView {
		t.Fatalf("active view = %v, want %v", m.activeView, invalidView)
	}
	line := statusViewLine(t, m)
	got := strings.TrimLeft(line, " ")
	if !strings.HasPrefix(got, programName) {
		t.Fatalf("status bar %q does not start with %q", line, programName)
	}
	if strings.HasPrefix(got, programName+breadcrumbSep) {
		t.Fatalf("status bar %q has a trailing breadcrumb separator", line)
	}
}

func TestStatusPlacesTitleLeftAndStatsRight(t *testing.T) {
	t.Parallel()
	ctx := testsize.Small(t)
	m := newStatusView(ctx)
	line := statusViewLine(t, m)
	want := programName + breadcrumbSep + tasksTitle
	if !strings.HasPrefix(strings.TrimLeft(line, " "), want) {
		t.Fatalf("status bar %q does not start with %q", line, want)
	}
	if !strings.HasSuffix(strings.TrimRight(line, " "), fakeStats) {
		t.Fatalf("status bar %q does not end with %q", line, fakeStats)
	}
}

func TestStatusMatchesWindowWidth(t *testing.T) {
	t.Parallel()
	ctx := testsize.Small(t)
	m := newStatusView(ctx)
	for _, width := range []int{1, 10, 80, 100, 120} {
		sized, _ := m.Update(tea.WindowSizeMsg{Width: width, Height: 24})
		got := lipgloss.Width(sized.View())
		if got != width {
			t.Fatalf("status bar width = %d, want %d", got, width)
		}
	}
}

func TestStatusDegradesOnNarrowWidth(t *testing.T) {
	t.Parallel()
	ctx := testsize.Small(t)
	m, _ := newStatusView(ctx).Update(tea.WindowSizeMsg{Width: 10, Height: 24})
	got := m.View()
	if lipgloss.Width(got) != 10 {
		t.Fatalf("status bar width = %d, want 10", lipgloss.Width(got))
	}
	if lipgloss.Height(got) != 2 {
		t.Fatalf("status bar height = %d, want 2", lipgloss.Height(got))
	}
	line := statusViewLine(t, m)
	content := strings.TrimSpace(line)
	if !strings.HasPrefix(fakeStats, content) {
		t.Fatalf("status bar %q, want a prefix of stats %q", line, fakeStats)
	}
}

func TestStatusTruncatesTitleBeforeStats(t *testing.T) {
	t.Parallel()
	ctx := testsize.Small(t)
	m := newStatusView(ctx)
	statsWidth := lipgloss.Width(m.theme.StatusStats.Inline(true).Render(fakeStats))
	frame := m.theme.StatusBar.GetHorizontalFrameSize()
	pad := strings.Repeat(" ", m.theme.StatusBar.GetPaddingLeft())

	t.Run("stats only", func(t *testing.T) {
		t.Parallel()
		sized, _ := m.Update(tea.WindowSizeMsg{Width: statsWidth + frame, Height: 24})
		line := statusViewLine(t, sized)
		want := pad + fakeStats + pad
		if line != want {
			t.Fatalf("status line = %q, want padded stats %q", line, want)
		}
	})

	t.Run("narrower than stats", func(t *testing.T) {
		t.Parallel()
		sized, _ := m.Update(tea.WindowSizeMsg{Width: 10, Height: 24})
		line := statusViewLine(t, sized)
		content := strings.TrimSpace(line)
		if !strings.HasPrefix(fakeStats, content) {
			t.Fatalf("status line = %q, want a prefix of stats %q", line, fakeStats)
		}
		if strings.HasPrefix(content, programName) {
			t.Fatalf("status line = %q starts with title, want truncated stats", line)
		}
	})

	t.Run("room for truncated title", func(t *testing.T) {
		t.Parallel()
		sized, _ := m.Update(tea.WindowSizeMsg{Width: statsWidth + 3 + frame, Height: 24})
		line := statusViewLine(t, sized)
		if !strings.HasPrefix(line, pad+programName[:3]) {
			t.Fatalf("status line = %q does not start with padded truncated title", line)
		}
		if !strings.HasSuffix(line, fakeStats+pad) {
			t.Fatalf("status line = %q does not end with padded stats %q", line, fakeStats)
		}
	})
}

func TestStatusUsesStyleGuideRoles(t *testing.T) {
	t.Parallel()
	ctx := testsize.Small(t)
	assertStatusRoles(t, newStatusView(ctx), styles.New(colorprofile.TrueColor))
}

func TestStatusAppliesUIThemeMsg(t *testing.T) {
	t.Parallel()
	ctx := testsize.Small(t)
	want := styles.New(colorprofile.ANSI)
	m := applyStatus(t, newStatusView(ctx), UIThemeMsg{Theme: want})
	assertStatusRoles(t, m, want)
}

func TestStatusDrawsBottomBorder(t *testing.T) {
	t.Parallel()
	ctx := testsize.Small(t)
	m := newStatusView(ctx)
	bar := visible(m.View())
	barLines := strings.Split(bar, "\n")
	if len(barLines) != 2 {
		t.Fatalf("status bar lines = %d, want 2", len(barLines))
	}
	rule := strings.Repeat(lipgloss.NormalBorder().Bottom, m.width)
	if barLines[1] != rule {
		t.Fatalf("status bar border = %q, want %q", barLines[1], rule)
	}
}

func TestStatusViewDoesNotMutateTheme(t *testing.T) {
	t.Parallel()
	ctx := testsize.Small(t)
	m := newStatusView(ctx)
	if m.theme.StatusBar.GetWidth() != 0 {
		t.Fatalf("StatusBar width = %d, want 0 before View()", m.theme.StatusBar.GetWidth())
	}
	_ = m.View()
	if m.theme.StatusBar.GetWidth() != 0 {
		t.Fatalf("StatusBar width = %d, want 0 after View()", m.theme.StatusBar.GetWidth())
	}
}

func TestStatusResizeStoresWidth(t *testing.T) {
	t.Parallel()
	ctx := testsize.Small(t)
	m := applyStatus(t, newStatusView(ctx), tea.WindowSizeMsg{Width: 100, Height: 20})
	if m.width != 100 {
		t.Fatalf("width = %d, want 100", m.width)
	}
}

func TestStatusUpdateAppliesAllMessages(t *testing.T) {
	t.Parallel()
	ctx := testsize.Small(t)
	m := applyStatus(t, newStatusView(ctx),
		tea.WindowSizeMsg{Width: 100, Height: 20},
		StatusViewActiveViewMsg{View: ViewTypeExample},
		UIThemeMsg{Theme: styles.New(colorprofile.ANSI)},
	)
	if m.width != 100 {
		t.Fatalf("width = %d, want 100", m.width)
	}
	if m.activeView != ViewTypeExample {
		t.Fatalf("active view = %v, want %v", m.activeView, ViewTypeExample)
	}
	assertStatusRoles(t, m, styles.New(colorprofile.ANSI))
}

func TestNewStatusDoesNotGetCounts(t *testing.T) {
	t.Parallel()
	ctx := testsize.Small(t)
	store := &fakeStore{nTasks: 3, nRunningAgents: 2}
	_ = NewStatusView(ctx, store, styles.New(colorprofile.TrueColor))
	if store.taskCountCalls != 0 {
		t.Fatalf("GetTaskCount() calls = %d, want 0", store.taskCountCalls)
	}
	if store.runningAgentCountCalls != 0 {
		t.Fatalf("GetRunningAgentCount() calls = %d, want 0", store.runningAgentCountCalls)
	}
}

func TestStatusFetchesStatsThroughCommand(t *testing.T) {
	t.Parallel()
	ctx := testsize.Small(t)
	store := &fakeStore{nTasks: 7, nRunningAgents: 4}
	m := NewStatusView(ctx, store, styles.New(colorprofile.TrueColor))
	m = applyStatus(t, m,
		tea.WindowSizeMsg{Width: 100, Height: 20},
		StatusViewActiveViewMsg{View: ViewTypeExample},
	)
	if store.taskCountCalls != 0 {
		t.Fatalf("GetTaskCount() calls before command = %d, want 0", store.taskCountCalls)
	}
	if store.runningAgentCountCalls != 0 {
		t.Fatalf("GetRunningAgentCount() calls before command = %d, want 0", store.runningAgentCountCalls)
	}
	m = applyStatus(t, m, m.statsCmd()())
	if store.taskCountCalls != 1 {
		t.Fatalf("GetTaskCount() calls = %d, want 1", store.taskCountCalls)
	}
	if store.runningAgentCountCalls != 1 {
		t.Fatalf("GetRunningAgentCount() calls = %d, want 1", store.runningAgentCountCalls)
	}
	want := "7 tasks, 4 agents running"
	line := statusViewLine(t, m)
	if !strings.Contains(line, want) {
		t.Fatalf("status bar %q does not contain stats %q", line, want)
	}
}

func TestStatusKeepsStatsOnTaskCountError(t *testing.T) {
	t.Parallel()
	ctx := testsize.Small(t)
	const previous = "7 tasks, 4 agents running"
	store := &fakeStore{nTasks: 7, nRunningAgents: 4}
	m := NewStatusView(ctx, store, styles.New(colorprofile.TrueColor))
	m = applyStatus(t, m, m.statsCmd()())
	store.taskCountErr = errors.New("task count failed")
	store.nTasks = 9
	store.nRunningAgents = 1
	m = applyStatus(t, m, m.statsCmd()())
	if store.taskCountCalls != 2 {
		t.Fatalf("GetTaskCount() calls = %d, want 2", store.taskCountCalls)
	}
	if store.runningAgentCountCalls != 1 {
		t.Fatalf("GetRunningAgentCount() calls = %d, want 1", store.runningAgentCountCalls)
	}
	if m.stats != previous {
		t.Fatalf("stats = %q, want previous value after GetTaskCount error", m.stats)
	}
	line := statusViewLine(t, m)
	if !strings.Contains(line, previous) {
		t.Fatalf("status bar %q does not contain previous stats", line)
	}
	if strings.Contains(line, "9 tasks") {
		t.Fatalf("status bar %q shows stats from a failed refresh", line)
	}
}

func TestStatusKeepsStatsOnRunningAgentCountError(t *testing.T) {
	t.Parallel()
	ctx := testsize.Small(t)
	const previous = "7 tasks, 4 agents running"
	store := &fakeStore{nTasks: 7, nRunningAgents: 4}
	m := NewStatusView(ctx, store, styles.New(colorprofile.TrueColor))
	m = applyStatus(t, m, m.statsCmd()())
	store.runningAgentCountErr = errors.New("agent count failed")
	store.nTasks = 9
	store.nRunningAgents = 1
	m = applyStatus(t, m, m.statsCmd()())
	if store.taskCountCalls != 2 {
		t.Fatalf("GetTaskCount() calls = %d, want 2", store.taskCountCalls)
	}
	if store.runningAgentCountCalls != 2 {
		t.Fatalf("GetRunningAgentCount() calls = %d, want 2", store.runningAgentCountCalls)
	}
	if m.stats != previous {
		t.Fatalf("stats = %q, want previous value after GetRunningAgentCount error", m.stats)
	}
	line := statusViewLine(t, m)
	if !strings.Contains(line, previous) {
		t.Fatalf("status bar %q does not contain previous stats", line)
	}
	if strings.Contains(line, "9 tasks") {
		t.Fatalf("status bar %q shows stats from a failed refresh", line)
	}
}

func TestStatusKeepsStatsOnCanceledContext(t *testing.T) {
	t.Parallel()
	ctx := testsize.Small(t)
	const previous = "7 tasks, 4 agents running"
	store := &fakeStore{nTasks: 7, nRunningAgents: 4}
	m := NewStatusView(ctx, store, styles.New(colorprofile.TrueColor))
	m = applyStatus(t, m, m.statsCmd()())
	canceled, cancel := context.WithCancel(ctx)
	cancel()
	m.ctx = canceled
	store.nTasks = 9
	store.nRunningAgents = 1
	m = applyStatus(t, m, m.statsCmd()())
	if m.stats != previous {
		t.Fatalf("stats = %q, want previous value after canceled context", m.stats)
	}
	line := statusViewLine(t, m)
	if !strings.Contains(line, previous) {
		t.Fatalf("status bar %q does not contain previous stats", line)
	}
	if strings.Contains(line, "9 tasks") {
		t.Fatalf("status bar %q shows stats from a canceled refresh", line)
	}
}

func TestStatusInitSchedulesRefresh(t *testing.T) {
	t.Parallel()
	ctx := testsize.Small(t)
	m := NewStatusView(ctx, &fakeStore{}, styles.New(colorprofile.TrueColor))
	if m.Init() == nil {
		t.Fatal("Init() cmd = nil, want stats refresh and tick")
	}
}

func TestStatusTickSchedulesRefresh(t *testing.T) {
	t.Parallel()
	ctx := testsize.Small(t)
	m := NewStatusView(ctx, &fakeStore{}, styles.New(colorprofile.TrueColor))
	_, cmd := m.Update(statusTickMsg{})
	if cmd == nil {
		t.Fatal("status tick cmd = nil, want stats refresh and next tick")
	}
}

func assertStatusRoles(t *testing.T, m StatusView, want styles.Theme) {
	t.Helper()
	p := want.Colors
	if m.theme.StatusTitle.GetForeground() != p.FG {
		t.Fatalf("status title foreground = %v, want fg %v", m.theme.StatusTitle.GetForeground(), p.FG)
	}
	if !m.theme.StatusTitle.GetBold() {
		t.Fatal("status title bold = false, want true")
	}
	if m.theme.StatusStats.GetForeground() != p.FGMuted {
		t.Fatalf("status stats foreground = %v, want fg-muted %v", m.theme.StatusStats.GetForeground(), p.FGMuted)
	}
	if _, ok := m.theme.StatusBar.GetBackground().(lipgloss.NoColor); !ok {
		t.Fatalf("status bar background = %v, want no color", m.theme.StatusBar.GetBackground())
	}
	if m.theme.StatusBar.GetBorderBottomForeground() != p.Border {
		t.Fatalf("status bar border = %v, want border %v", m.theme.StatusBar.GetBorderBottomForeground(), p.Border)
	}
}
