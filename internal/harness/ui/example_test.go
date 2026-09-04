package ui

import (
	"slices"
	"testing"

	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"github.com/charlieparkes/go-testsize"
	"github.com/charlieparkes/harness/internal/harness/ui/styles"
	"github.com/charmbracelet/colorprofile"
)

func newExampleView() ExampleView {
	return NewExampleView(styles.New(colorprofile.TrueColor))
}

func applyExample(t *testing.T, m ExampleView, msgs ...tea.Msg) ExampleView {
	t.Helper()
	for _, msg := range msgs {
		m, _ = m.Update(msg)
	}
	return m
}

func updateExample(t *testing.T, m ExampleView, msgs ...tea.Msg) (ExampleView, tea.Cmd) {
	t.Helper()
	cmds := make([]tea.Cmd, 0, len(msgs))
	for _, msg := range msgs {
		var cmd tea.Cmd
		m, cmd = m.Update(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}
	return m, tea.Batch(cmds...)
}

func TestExampleName(t *testing.T) {
	t.Parallel()
	testsize.Small(t)
	got := ExampleView{}.Name()
	if got != "Example" {
		t.Fatalf("Name() = %q, want Example", got)
	}
}

func TestExampleUsesMiniDot(t *testing.T) {
	t.Parallel()
	testsize.Small(t)
	m := newExampleView()
	if !slices.Equal(m.spinner.Spinner.Frames, spinner.MiniDot.Frames) {
		t.Fatalf("spinner frames = %q, want MiniDot %q", m.spinner.Spinner.Frames, spinner.MiniDot.Frames)
	}
	if m.spinner.Spinner.FPS != spinner.MiniDot.FPS {
		t.Fatalf("spinner FPS = %s, want %s", m.spinner.Spinner.FPS, spinner.MiniDot.FPS)
	}
}

func TestExampleUsesMintSpinner(t *testing.T) {
	t.Parallel()
	testsize.Small(t)
	assertSpinnerStyle(t, newExampleView().spinner, styles.New(colorprofile.TrueColor))
}

func TestExampleAppliesUIThemeMsg(t *testing.T) {
	t.Parallel()
	testsize.Small(t)
	want := styles.New(colorprofile.ANSI)
	m := applyExample(t, newExampleView(), UIThemeMsg{Theme: want})
	assertSpinnerStyle(t, m.spinner, want)
}

func TestExampleInitStartsTick(t *testing.T) {
	t.Parallel()
	testsize.Small(t)
	m := newExampleView()
	if m.Init() == nil {
		t.Fatal("Init() cmd = nil, want spinner tick")
	}
}

func TestExampleTickAdvancesFrame(t *testing.T) {
	t.Parallel()
	testsize.Small(t)
	m := newExampleView()
	first := visible(m.View())
	if first != spinner.MiniDot.Frames[0] {
		t.Fatalf("initial view = %q, want %q", first, spinner.MiniDot.Frames[0])
	}

	var cmd tea.Cmd
	m, cmd = updateExample(t, m, m.spinner.Tick())
	if cmd == nil {
		t.Fatal("first tick cmd = nil, want next spinner tick")
	}
	if visible(m.View()) != spinner.MiniDot.Frames[1] {
		t.Fatalf("after first tick view = %q, want %q", m.View(), spinner.MiniDot.Frames[1])
	}

	m, cmd = updateExample(t, m, m.spinner.Tick())
	if cmd == nil {
		t.Fatal("second tick cmd = nil, want next spinner tick")
	}
	if visible(m.View()) != spinner.MiniDot.Frames[2] {
		t.Fatalf("after second tick view = %q, want %q", m.View(), spinner.MiniDot.Frames[2])
	}
}

func TestExampleResizeStoresSize(t *testing.T) {
	t.Parallel()
	testsize.Small(t)
	m := applyExample(t, newExampleView(), tea.WindowSizeMsg{Width: 100, Height: 20})
	if m.width != 100 {
		t.Fatalf("width = %d, want 100", m.width)
	}
	if m.height != 20 {
		t.Fatalf("height = %d, want 20", m.height)
	}
}

func TestExampleUpdateAppliesAllMessages(t *testing.T) {
	t.Parallel()
	testsize.Small(t)
	m := newExampleView()
	tick := m.spinner.Tick()
	m, cmd := updateExample(t, m, tea.WindowSizeMsg{Width: 100, Height: 20}, tick)
	if m.width != 100 {
		t.Fatalf("width = %d, want 100", m.width)
	}
	if m.height != 20 {
		t.Fatalf("height = %d, want 20", m.height)
	}
	if visible(m.View()) != spinner.MiniDot.Frames[1] {
		t.Fatalf("view = %q, want %q", m.View(), spinner.MiniDot.Frames[1])
	}
	if cmd == nil {
		t.Fatal("Update() cmd = nil, want next spinner tick")
	}
}

func assertSpinnerStyle(t *testing.T, s spinner.Model, want styles.Theme) {
	t.Helper()
	if s.Style.GetForeground() != want.Colors.Success {
		t.Fatalf("spinner foreground = %v, want success %v", s.Style.GetForeground(), want.Colors.Success)
	}
}
