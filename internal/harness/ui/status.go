package ui

import (
	"context"
	"fmt"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charlieparkes/harness/internal/harness/ui/styles"
)

const (
	programName           = "Harness"
	breadcrumbSep         = " / "
	statusRefreshInterval = 5 * time.Second
)

// StatusStore is the persistence needed by the Status view.
type StatusStore interface {
	GetTaskCount(ctx context.Context) (int64, error)
	GetRunningAgentCount(ctx context.Context) (int64, error)
}

// StatusViewActiveViewMsg tells StatusView which child view is active.
type StatusViewActiveViewMsg struct {
	View ViewType
}

type statusTickMsg struct{}

type statusStatsMsg struct {
	nTasks         int64
	nRunningAgents int64
	err            error
}

// StatusView is the application status bar.
type StatusView struct {
	ctx        context.Context
	store      StatusStore
	width      int
	theme      styles.Theme
	activeView ViewType
	stats      string
}

func NewStatusView(ctx context.Context, store StatusStore, theme styles.Theme) StatusView {
	return StatusView{
		ctx:   ctx,
		store: store,
		width: defaultWidth,
		theme: theme,
	}
}

func (m StatusView) Init() tea.Cmd {
	return m.refresh()
}

func (m StatusView) Update(msg tea.Msg) (StatusView, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
	case StatusViewActiveViewMsg:
		m.activeView = msg.View
	case UIThemeMsg:
		m.theme = msg.Theme
	case statusTickMsg:
		return m, m.refresh()
	case statusStatsMsg:
		if msg.err == nil {
			m.stats = formatStats(msg.nTasks, msg.nRunningAgents)
		}
	}
	return m, nil
}

func (m StatusView) refresh() tea.Cmd {
	return tea.Batch(m.statsCmd(), m.tickCmd())
}

func (m StatusView) statsCmd() tea.Cmd {
	return func() tea.Msg {
		nTasks, err := m.store.GetTaskCount(m.ctx)
		if err != nil {
			return statusStatsMsg{err: err}
		}
		nRunningAgents, err := m.store.GetRunningAgentCount(m.ctx)
		if err != nil {
			return statusStatsMsg{err: err}
		}
		return statusStatsMsg{
			nTasks:         nTasks,
			nRunningAgents: nRunningAgents,
		}
	}
}

func (m StatusView) tickCmd() tea.Cmd {
	return tea.Tick(statusRefreshInterval, func(time.Time) tea.Msg {
		return statusTickMsg{}
	})
}

func (m StatusView) View() string {
	width := max(m.width, 1)
	inner := max(width-m.theme.StatusBar.GetHorizontalFrameSize(), 1)
	right := m.theme.StatusStats.
		Inline(true).
		MaxWidth(inner).
		Render(m.stats)
	remain := max(inner-lipgloss.Width(right), 0)
	row := right
	if remain > 0 {
		content := m.theme.StatusTitle.
			Inline(true).
			Render(programName)
		if name := m.activeView.Name(); name != "" {
			content += m.theme.StatusBreadcrumbs.
				Inline(true).
				Render(breadcrumbSep + name)
		}
		left := lipgloss.NewStyle().
			Width(remain).
			MaxWidth(remain).
			Inline(true).
			Render(content)
		row = lipgloss.JoinHorizontal(lipgloss.Top, left, right)
	}
	return m.theme.StatusBar.
		Width(width).
		MaxWidth(width).
		Render(row)
}

func formatStats(nTasks, nRunningAgents int64) string {
	return fmt.Sprintf("%d tasks, %d agents running", nTasks, nRunningAgents)
}
