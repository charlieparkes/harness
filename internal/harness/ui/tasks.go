package ui

import (
	"context"
	"io"
	"strings"
	"time"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charlieparkes/harness/internal/harness/model"
	"github.com/charlieparkes/harness/internal/harness/ui/styles"
	"github.com/charlieparkes/harness/internal/shared"
)

const (
	refreshInterval = 5 * time.Second
	tasksTitle      = "Tasks"
)

// TasksStore is the persistence needed by the Tasks view.
type TasksStore interface {
	ListTasks(ctx context.Context) ([]model.Task, error)
}

type tasksTickMsg struct{}

type listTasksMsg struct {
	tasks []model.Task
	err   error
}

type taskItem struct {
	task         model.Task
	theme        styles.Theme
	selected     bool
	runningFrame string
}

var _ list.DefaultItem = taskItem{}

func (i taskItem) Title() string {
	return i.task.Title
}

func (i taskItem) Description() string {
	return i.descriptionAt(time.Now())
}

func (i taskItem) descriptionAt(now time.Time) string {
	parts := make([]string, 0, 2)
	if i.task.Status == model.TaskStatusRunning {
		parts = append(parts, i.runningFrame+" "+i.theme.StatusLabel(i.selected, i.task.Status).Render(i.task.Status.String()))
	} else if label := i.task.Status.Text(); label != "" {
		parts = append(parts, i.theme.StatusLabel(i.selected, i.task.Status).Render(label))
	}
	if rel := shared.RelativeTime(i.task.CreatedAt, now); rel != "" {
		parts = append(parts, lipgloss.NewStyle().
			Foreground(i.theme.Colors.FGMuted).
			Inline(true).
			Render(rel))
	}
	return strings.Join(parts, " ")
}

func (i taskItem) FilterValue() string {
	parts := []string{i.task.Title, i.task.ID}
	if status := i.task.Status.String(); status != "" {
		parts = append(parts, status)
	}
	return strings.Join(parts, " ")
}

func taskItems(tasks []model.Task) []list.Item {
	items := make([]list.Item, len(tasks))
	for i, task := range tasks {
		items[i] = taskItem{task: task}
	}
	return items
}

type taskDelegate struct {
	list.DefaultDelegate
	theme        styles.Theme
	runningFrame string
}

var _ list.ItemDelegate = taskDelegate{}

func (d taskDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	if t, ok := item.(taskItem); ok {
		d.Styles = d.theme.TaskItems(t.task.Status)
		item = taskItem{
			task:         t.task,
			theme:        d.theme,
			selected:     index == m.Index() && m.FilterState() != list.Filtering,
			runningFrame: d.runningFrame,
		}
	}
	d.DefaultDelegate.Render(w, m, index, item)
}

// TasksView is a selectable list of tasks.
type TasksView struct {
	ctx      context.Context
	store    TasksStore
	list     list.Model
	delegate taskDelegate
	spinner  spinner.Model
	theme    styles.Theme
	quitting bool
}

var _ UIView = TasksView{}

func NewTasksView(ctx context.Context, store TasksStore, theme styles.Theme) TasksView {
	delegate := taskDelegate{DefaultDelegate: list.NewDefaultDelegate()}
	l := list.New(nil, delegate, defaultWidth, defaultHeight)
	l.SetShowTitle(false)
	l.SetShowStatusBar(false)
	l.SetStatusBarItemName("task", "tasks")

	km := list.DefaultKeyMap()
	km.Quit = key.NewBinding(
		key.WithKeys("q"),
		key.WithHelp("q", "quit"),
	)
	l.KeyMap = km

	m := TasksView{
		ctx:      ctx,
		store:    store,
		list:     l,
		delegate: delegate,
		spinner: spinner.New(
			spinner.WithSpinner(spinner.MiniDot),
			spinner.WithStyle(theme.Spinner),
		),
		theme: theme,
	}
	return m.setDelegateStyles()
}

func (m TasksView) Name() string {
	return tasksTitle
}

func (m TasksView) Init() tea.Cmd {
	return tea.Batch(m.show(), m.spinner.Tick)
}

func (m TasksView) Update(msg tea.Msg) (TasksView, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		}
	case tea.InterruptMsg:
		m.quitting = true
		return m, tea.Quit
	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height)
		return m, nil
	case UIThemeMsg:
		m.theme = msg.Theme
		m.spinner.Style = msg.Theme.Spinner
		return m.setDelegateStyles(), nil
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		if cmd == nil {
			return m, nil
		}
		return m.setDelegateStyles(), cmd
	case tasksTickMsg:
		return m, m.show()
	case listTasksMsg:
		if m.quitting || msg.err != nil {
			return m, nil
		}
		var cmd tea.Cmd
		m, cmd = m.setTasks(msg.tasks)
		return m, cmd
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m TasksView) View() string {
	if m.quitting {
		return ""
	}

	return m.list.View()
}

func (m TasksView) show() tea.Cmd {
	if m.quitting {
		return nil
	}
	return tea.Batch(m.listTasksCmd(), m.tickCmd())
}

func (m TasksView) tickCmd() tea.Cmd {
	return tea.Tick(refreshInterval, func(time.Time) tea.Msg {
		return tasksTickMsg{}
	})
}

func (m TasksView) listTasksCmd() tea.Cmd {
	return func() tea.Msg {
		tasks, err := m.store.ListTasks(m.ctx)
		return listTasksMsg{tasks: tasks, err: err}
	}
}

func (m TasksView) setDelegateStyles() TasksView {
	m.delegate.theme = m.theme
	m.delegate.Styles = m.theme.TaskItems(model.TaskStatusReady)
	m.delegate.runningFrame = m.spinner.View()
	m.list.SetDelegate(m.delegate)
	return m
}

func (m TasksView) setTasks(tasks []model.Task) (TasksView, tea.Cmd) {
	var selectedID string
	if item, ok := m.list.SelectedItem().(taskItem); ok {
		selectedID = item.task.ID
	}

	cmd := m.list.SetItems(taskItems(tasks))
	m.list.Select(0)
	for i, item := range m.list.Items() {
		if t, ok := item.(taskItem); ok && t.task.ID == selectedID {
			m.list.Select(i)
			break
		}
	}
	return m, cmd
}
