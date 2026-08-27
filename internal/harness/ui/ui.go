package ui

import (
	"context"
	"os"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charlieparkes/harness/internal/harness/ui/styles"
	"github.com/charmbracelet/colorprofile"
)

const (
	defaultWidth  = 80
	defaultHeight = 24
)

// Store is the persistence needed by the application shell.
type Store interface {
	TasksStore
	StatusStore
}

// UIView is a renderable child view of the application shell.
//
//nolint:revive // Package is ui; UIView is the child-view contract for Model.
type UIView interface {
	Name() string
	View() string
}

// UIThemeMsg tells a child view to apply a theme.
//
//nolint:revive // Package is ui; UIThemeMsg is the shared theme-update message.
type UIThemeMsg struct {
	Theme styles.Theme
}

// Model is the application shell.
type Model struct {
	width      int
	height     int
	theme      styles.Theme
	activeView ViewType
	status     StatusView
	tasks      TasksView
	example    ExampleView
	quitting   bool
}

func New(ctx context.Context, store Store) Model {
	theme := styles.New(colorprofile.Detect(os.Stderr, os.Environ()))
	m := Model{
		width:   defaultWidth,
		height:  defaultHeight,
		theme:   theme,
		status:  NewStatusView(ctx, store, theme),
		tasks:   NewTasksView(ctx, store, theme),
		example: NewExampleView(theme),
	}
	m = m.setActiveView(ViewTypeTasks)
	m, _ = m.resizeChildren()
	return m
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.status.Init(), m.tasks.Init(), m.example.Init())
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			m.tasks, _ = m.tasks.Update(msg)
			return m, tea.Quit
		}
		return m.updateActiveChild(msg)
	case tea.ColorProfileMsg:
		m.theme = styles.New(msg.Profile)
		return m.updateAll(UIThemeMsg{Theme: m.theme})
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m.resizeChildren()
	case tea.InterruptMsg:
		m.quitting = true
		m.tasks, _ = m.tasks.Update(msg)
		return m, tea.Quit
	}
	return m.updateAll(msg)
}

func (m Model) setActiveView(view ViewType) Model {
	m.activeView = view
	m.status, _ = m.status.Update(StatusViewActiveViewMsg{View: m.activeView})
	return m
}

func (m Model) updateActiveChild(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	switch m.activeView {
	case ViewTypeTasks:
		m.tasks, cmd = m.tasks.Update(msg)
	case ViewTypeExample:
		m.example, cmd = m.example.Update(msg)
	}
	return m, cmd
}

func (m Model) updateAll(msg tea.Msg) (Model, tea.Cmd) {
	var statusCmd, tasksCmd, exampleCmd tea.Cmd
	m.status, statusCmd = m.status.Update(msg)
	m.tasks, tasksCmd = m.tasks.Update(msg)
	m.example, exampleCmd = m.example.Update(msg)
	return m, tea.Batch(statusCmd, tasksCmd, exampleCmd)
}

func (m Model) resizeChildren() (Model, tea.Cmd) {
	var statusCmd, tasksCmd, exampleCmd tea.Cmd
	m.status, statusCmd = m.status.Update(tea.WindowSizeMsg{Width: m.width, Height: m.height})
	w, h := m.childWindowSize()
	size := tea.WindowSizeMsg{Width: w, Height: h}
	m.tasks, tasksCmd = m.tasks.Update(size)
	m.example, exampleCmd = m.example.Update(size)
	return m, tea.Batch(statusCmd, tasksCmd, exampleCmd)
}

func (m Model) childWindowSize() (int, int) {
	w := max(m.width, 1)
	h := max(m.height-lipgloss.Height(m.status.View()), 1)
	return w, h
}

func (m Model) View() tea.View {
	v := tea.NewView(m.render())
	v.AltScreen = true
	return v
}

func (m Model) render() string {
	if m.quitting {
		return ""
	}
	return lipgloss.JoinVertical(lipgloss.Left, m.status.View(), m.activeChild())
}

func (m Model) activeChild() string {
	if v := m.activeChildView(); v != nil {
		return v.View()
	}
	return ""
}

func (m Model) activeChildView() UIView {
	switch m.activeView {
	case ViewTypeTasks:
		return m.tasks
	case ViewTypeExample:
		return m.example
	default:
		return nil
	}
}
