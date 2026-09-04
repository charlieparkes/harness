package ui

import (
	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"github.com/charlieparkes/harness/internal/harness/ui/styles"
)

const exampleTitle = "Example"

// ExampleView is a MiniDot spinner shown on the ExampleView tab.
type ExampleView struct {
	spinner spinner.Model
	width   int
	height  int
}

var _ UIView = ExampleView{}

func NewExampleView(theme styles.Theme) ExampleView {
	return ExampleView{
		spinner: spinner.New(
			spinner.WithSpinner(spinner.MiniDot),
			spinner.WithStyle(theme.Spinner),
		),
		width:  defaultWidth,
		height: defaultHeight,
	}
}

func (m ExampleView) Name() string {
	return exampleTitle
}

func (m ExampleView) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m ExampleView) Update(msg tea.Msg) (ExampleView, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case UIThemeMsg:
		m.spinner.Style = msg.Theme.Spinner
		return m, nil
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m ExampleView) View() string {
	return m.spinner.View()
}
