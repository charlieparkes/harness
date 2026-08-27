package styles

import (
	"image/color"

	"charm.land/bubbles/v2/list"
	"charm.land/lipgloss/v2"
	"github.com/charlieparkes/harness/internal/harness/model"
	"github.com/charmbracelet/colorprofile"
)

const (
	commonPad      = 2
	selectedMarker = "❯"
)

// Theme is the Lip Gloss styles for current TUI surfaces.
type Theme struct {
	Colors            Palette
	Spinner           lipgloss.Style
	StatusBar         lipgloss.Style
	StatusTitle       lipgloss.Style
	StatusBreadcrumbs lipgloss.Style
	StatusStats       lipgloss.Style
}

// New returns a theme resolved for profile.
func New(profile colorprofile.Profile) Theme {
	colors := resolve(profile)
	return Theme{
		Colors: colors,
		StatusBar: lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, false, true, false).
			BorderForeground(colors.Border).
			Padding(0, commonPad),
		StatusTitle: lipgloss.NewStyle().
			Foreground(colors.FG).
			Bold(true),
		StatusBreadcrumbs: lipgloss.NewStyle().
			Foreground(colors.FGMuted),
		StatusStats: lipgloss.NewStyle().
			Foreground(colors.FGMuted),
		Spinner: lipgloss.NewStyle().
			Foreground(colors.Success),
	}
}

// TaskItems returns list item styles for a task status. The selected title
// uses a "❯" marker. Title foreground follows status.
func (t Theme) TaskItems(status model.TaskStatus) list.DefaultItemStyles {
	padded := lipgloss.NewStyle().Padding(0, 0, 0, commonPad)

	title := padded.
		Foreground(t.titleColor(status))
	selectedTitle := padded.
		Border(lipgloss.Border{Left: selectedMarker}, false, false, false, true).
		BorderForeground(t.Colors.FG).
		Padding(0, 0, 0, commonPad-1).
		Foreground(t.Colors.FG)
	dimmedTitle := padded.
		Foreground(t.Colors.FGMuted)

	desc := padded.
		Foreground(t.Colors.FGMuted)
	selectedDesc := padded.
		Foreground(t.Colors.FGMuted)
	dimmedDesc := padded.
		Foreground(t.Colors.FGMuted)

	filterMatch := lipgloss.NewStyle().
		Foreground(t.Colors.Accent)

	return list.DefaultItemStyles{
		NormalTitle:   title,
		SelectedTitle: selectedTitle,
		DimmedTitle:   dimmedTitle,

		NormalDesc:   desc,
		SelectedDesc: selectedDesc,
		DimmedDesc:   dimmedDesc,

		FilterMatch: filterMatch,
	}
}

// StatusLabel returns the style for a task description status name.
func (t Theme) StatusLabel(selected bool, status model.TaskStatus) lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(t.StatusColor(selected, status)).
		Inline(true)
}

func (t Theme) titleColor(status model.TaskStatus) color.Color {
	switch status {
	case model.TaskStatusReady, model.TaskStatusRunning, model.TaskStatusPendingFeedback, model.TaskStatusFailed:
		return t.Colors.FGMuted
	case model.TaskStatusUnspecified, model.TaskStatusCanceled, model.TaskStatusCompleted:
		return t.Colors.FGMuted
	default:
		return t.Colors.FGMuted
	}
}

func (t Theme) StatusColor(selected bool, status model.TaskStatus) color.Color {
	switch status {
	case model.TaskStatusReady:
		return t.Colors.FG
	case model.TaskStatusRunning:
		if selected {
			return t.Colors.Success
		}
		return t.Colors.SuccessMuted
	case model.TaskStatusPendingFeedback:
		if selected {
			return t.Colors.Warning
		}
		return t.Colors.WarningMuted
	case model.TaskStatusFailed:
		if selected {
			return t.Colors.Danger
		}
		return t.Colors.DangerMuted
	case model.TaskStatusCanceled, model.TaskStatusCompleted:
		return t.Colors.FGMuted
	case model.TaskStatusUnspecified:
		return t.Colors.Danger
	default:
		return t.Colors.FG
	}
}
