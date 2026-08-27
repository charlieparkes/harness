package styles

import (
	"image/color"
	"regexp"
	"strings"
	"testing"

	"charm.land/lipgloss/v2"
	"github.com/charlieparkes/go-testsize"
	"github.com/charlieparkes/harness/internal/harness/model"
	"github.com/charmbracelet/colorprofile"
)

var ansiEscape = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func visible(s string) string {
	return ansiEscape.ReplaceAllString(s, "")
}

func TestThemeCanvasAndStatusRoles(t *testing.T) {
	t.Parallel()
	testsize.Small(t)

	profiles := []colorprofile.Profile{
		colorprofile.TrueColor,
		colorprofile.ANSI256,
		colorprofile.ANSI,
	}
	for _, profile := range profiles {
		t.Run(profile.String(), func(t *testing.T) {
			t.Parallel()
			theme := New(profile)
			p := theme.Colors
			if theme.StatusTitle.GetForeground() != p.FG {
				t.Fatalf("status title foreground = %v, want fg %v", theme.StatusTitle.GetForeground(), p.FG)
			}
			if !theme.StatusTitle.GetBold() {
				t.Fatal("status title bold = false, want true")
			}
			if theme.StatusBreadcrumbs.GetForeground() != p.FGMuted {
				t.Fatalf("status breadcrumbs foreground = %v, want fg-muted %v", theme.StatusBreadcrumbs.GetForeground(), p.FGMuted)
			}
			if theme.StatusStats.GetForeground() != p.FGMuted {
				t.Fatalf("status stats foreground = %v, want fg-muted %v", theme.StatusStats.GetForeground(), p.FGMuted)
			}
			assertNoBackground(t, theme.StatusBar, "status bar")
			if !theme.StatusBar.GetBorderBottom() {
				t.Fatal("status bar border bottom = false, want true")
			}
			if theme.StatusBar.GetBorderTop() || theme.StatusBar.GetBorderLeft() || theme.StatusBar.GetBorderRight() {
				t.Fatal("status bar has extra borders, want bottom only")
			}
			if theme.StatusBar.GetBorderBottomForeground() != p.Border {
				t.Fatalf("status bar border = %v, want border %v", theme.StatusBar.GetBorderBottomForeground(), p.Border)
			}
			if theme.StatusBar.GetPaddingLeft() != commonPad {
				t.Fatalf("status bar padding left = %d, want %d", theme.StatusBar.GetPaddingLeft(), commonPad)
			}
			if theme.StatusBar.GetPaddingRight() != commonPad {
				t.Fatalf("status bar padding right = %d, want %d", theme.StatusBar.GetPaddingRight(), commonPad)
			}
			if theme.Spinner.GetForeground() != p.Success {
				t.Fatalf("spinner foreground = %v, want success %v", theme.Spinner.GetForeground(), p.Success)
			}
		})
	}
}

func TestTaskItemRoles(t *testing.T) {
	t.Parallel()
	testsize.Small(t)

	theme := New(colorprofile.TrueColor)
	p := theme.Colors
	items := theme.TaskItems(model.TaskStatusReady)

	if items.NormalTitle.GetForeground() != p.FGMuted {
		t.Fatalf("normal title foreground = %v, want fg-muted %v", items.NormalTitle.GetForeground(), p.FGMuted)
	}
	assertNoBackground(t, items.NormalTitle, "normal title")
	if items.NormalDesc.GetForeground() != p.FGMuted {
		t.Fatalf("normal description foreground = %v, want fg-muted %v", items.NormalDesc.GetForeground(), p.FGMuted)
	}
	assertNoBackground(t, items.NormalDesc, "normal description")
	if items.DimmedTitle.GetForeground() != p.FGMuted {
		t.Fatalf("dimmed title foreground = %v, want fg-muted %v", items.DimmedTitle.GetForeground(), p.FGMuted)
	}
	assertNoBackground(t, items.DimmedTitle, "dimmed title")
	if items.DimmedDesc.GetForeground() != p.FGMuted {
		t.Fatalf("dimmed description foreground = %v, want fg-muted %v", items.DimmedDesc.GetForeground(), p.FGMuted)
	}
	assertNoBackground(t, items.DimmedDesc, "dimmed description")
	if items.FilterMatch.GetForeground() != p.Accent {
		t.Fatalf("filter match foreground = %v, want accent %v", items.FilterMatch.GetForeground(), p.Accent)
	}
	assertNoBackground(t, items.FilterMatch, "filter match")
	if items.SelectedTitle.GetForeground() != p.FG {
		t.Fatalf("selected title foreground = %v, want fg %v", items.SelectedTitle.GetForeground(), p.FG)
	}
	if items.SelectedTitle.GetBold() {
		t.Fatal("selected title bold = true, want false")
	}
	if items.NormalTitle.GetBold() {
		t.Fatal("normal title bold = true, want false")
	}
	if hasBoldSGR(items.SelectedTitle.Render("Item")) {
		t.Fatal("selected title render has bold SGR")
	}
	if hasBoldSGR(items.NormalTitle.Render("Item")) {
		t.Fatal("normal title render has bold SGR")
	}
	if items.SelectedDesc.GetForeground() != p.FGMuted {
		t.Fatalf("selected description foreground = %v, want fg-muted %v", items.SelectedDesc.GetForeground(), p.FGMuted)
	}
	assertNoBackground(t, items.SelectedTitle, "selected title")
	assertNoBackground(t, items.SelectedDesc, "selected description")
	assertSelectedTitleMarker(t, items.SelectedTitle, p, "selected title")
	assertPaddedNoBorder(t, items.SelectedDesc, "selected description")
	assertPaddedNoBorder(t, items.NormalTitle, "normal title")
	assertPaddedNoBorder(t, items.NormalDesc, "normal description")
	assertPaddedNoBorder(t, items.DimmedTitle, "dimmed title")
	assertPaddedNoBorder(t, items.DimmedDesc, "dimmed description")
	if items.NormalTitle.GetWidth() != 0 {
		t.Fatalf("normal title width = %d, want 0", items.NormalTitle.GetWidth())
	}
	if items.SelectedTitle.GetWidth() != 0 {
		t.Fatalf("selected title width = %d, want 0", items.SelectedTitle.GetWidth())
	}
}

func TestSelectedTitleAlignsWithNormal(t *testing.T) {
	t.Parallel()
	testsize.Small(t)

	theme := New(colorprofile.TrueColor)
	items := theme.TaskItems(model.TaskStatusReady)
	const title = "Title"
	normal := visible(items.NormalTitle.Render(title))
	selected := visible(items.SelectedTitle.Render(title))
	if !strings.HasPrefix(selected, selectedMarker) {
		t.Fatalf("selected title %q does not start with %q", selected, selectedMarker)
	}
	if strings.Contains(selected, lipgloss.NormalBorder().Left) {
		t.Fatalf("selected title %q still has a vertical rule", selected)
	}
	nBefore, _, nFound := strings.Cut(normal, title)
	sBefore, _, sFound := strings.Cut(selected, title)
	if !nFound || !sFound {
		t.Fatalf("title missing: normal %q selected %q", normal, selected)
	}
	if lipgloss.Width(nBefore) != lipgloss.Width(sBefore) {
		t.Fatalf("title column: normal %q selected %q", normal, selected)
	}
}

func TestSelectedDescriptionAlignsWithTitle(t *testing.T) {
	t.Parallel()
	testsize.Small(t)

	theme := New(colorprofile.TrueColor)
	items := theme.TaskItems(model.TaskStatusReady)
	const (
		title = "Title"
		desc  = "Ready 1d ago"
	)
	selectedTitle := visible(items.SelectedTitle.Render(title))
	selectedDesc := visible(items.SelectedDesc.Render(desc))
	normalDesc := visible(items.NormalDesc.Render(desc))
	if strings.Contains(selectedDesc, selectedMarker) {
		t.Fatalf("selected description %q has marker %q", selectedDesc, selectedMarker)
	}
	if strings.Contains(selectedDesc, lipgloss.NormalBorder().Left) {
		t.Fatalf("selected description %q still has a vertical rule", selectedDesc)
	}
	sBefore, _, sFound := strings.Cut(selectedTitle, title)
	dBefore, _, dFound := strings.Cut(selectedDesc, desc)
	nBefore, _, nFound := strings.Cut(normalDesc, desc)
	if !sFound || !dFound || !nFound {
		t.Fatalf("text missing: title %q selected desc %q normal desc %q", selectedTitle, selectedDesc, normalDesc)
	}
	if lipgloss.Width(sBefore) != lipgloss.Width(dBefore) {
		t.Fatalf("selected description column: title prefix %q desc prefix %q", sBefore, dBefore)
	}
	if lipgloss.Width(dBefore) != lipgloss.Width(nBefore) {
		t.Fatalf("description column: selected %q normal %q", selectedDesc, normalDesc)
	}
}

func TestDerivedStylesAreIndependent(t *testing.T) {
	t.Parallel()
	testsize.Small(t)

	theme := New(colorprofile.TrueColor)
	items := theme.TaskItems(model.TaskStatusReady)
	mutated := items.SelectedTitle.BorderForeground(theme.Colors.Border)
	if items.SelectedTitle.GetBorderLeftForeground() != theme.Colors.FG {
		t.Fatalf("SelectedTitle border = %v, want fg after deriving a copy", items.SelectedTitle.GetBorderLeftForeground())
	}
	if mutated.GetBorderLeftForeground() != theme.Colors.Border {
		t.Fatalf("derived border = %v, want border %v", mutated.GetBorderLeftForeground(), theme.Colors.Border)
	}
	if mutated.GetPaddingLeft() != items.SelectedTitle.GetPaddingLeft() {
		t.Fatalf("derived padding left = %d, want %d", mutated.GetPaddingLeft(), items.SelectedTitle.GetPaddingLeft())
	}
	if mutated.GetBorderLeft() != items.SelectedTitle.GetBorderLeft() {
		t.Fatalf("derived border left = %v, want %v", mutated.GetBorderLeft(), items.SelectedTitle.GetBorderLeft())
	}

	tinted := theme.StatusTitle.Foreground(theme.Colors.Accent)
	if theme.StatusTitle.GetForeground() != theme.Colors.FG {
		t.Fatalf("StatusTitle foreground = %v, want fg after deriving a tinted copy", theme.StatusTitle.GetForeground())
	}
	if tinted.GetForeground() != theme.Colors.Accent {
		t.Fatalf("derived title foreground = %v, want accent %v", tinted.GetForeground(), theme.Colors.Accent)
	}
	if !tinted.GetBold() {
		t.Fatal("derived title dropped bold")
	}

	if theme.StatusBar.GetWidth() != 0 {
		t.Fatalf("StatusBar width = %d, want 0", theme.StatusBar.GetWidth())
	}
}

func TestTaskItemTitleColorsByStatus(t *testing.T) {
	t.Parallel()
	testsize.Small(t)

	theme := New(colorprofile.TrueColor)
	p := theme.Colors
	tests := []struct {
		status    model.TaskStatus
		wantTitle color.Color
	}{
		{model.TaskStatusReady, p.FGMuted},
		{model.TaskStatusRunning, p.FGMuted},
		{model.TaskStatusPendingFeedback, p.FGMuted},
		{model.TaskStatusCompleted, p.FGMuted},
		{model.TaskStatusCanceled, p.FGMuted},
		{model.TaskStatusFailed, p.FGMuted},
		{model.TaskStatus("unknown"), p.FGMuted},
	}
	for _, tc := range tests {
		t.Run(tc.status.String(), func(t *testing.T) {
			t.Parallel()
			items := theme.TaskItems(tc.status)
			if items.NormalTitle.GetForeground() != tc.wantTitle {
				t.Fatalf("normal title foreground = %v, want %v", items.NormalTitle.GetForeground(), tc.wantTitle)
			}
			if items.SelectedTitle.GetForeground() != p.FG {
				t.Fatalf("selected title foreground = %v, want fg %v", items.SelectedTitle.GetForeground(), p.FG)
			}
			if items.SelectedTitle.GetBold() {
				t.Fatal("selected title bold = true, want false")
			}
			if items.NormalDesc.GetForeground() != p.FGMuted {
				t.Fatalf("normal description foreground = %v, want fg-muted %v", items.NormalDesc.GetForeground(), p.FGMuted)
			}
			if items.SelectedDesc.GetForeground() != p.FGMuted {
				t.Fatalf("selected description foreground = %v, want fg-muted %v", items.SelectedDesc.GetForeground(), p.FGMuted)
			}
			assertSelectedTitleMarker(t, items.SelectedTitle, p, "selected title")
			assertPaddedNoBorder(t, items.SelectedDesc, "selected description")
		})
	}
}

func TestStatusLabelColors(t *testing.T) {
	t.Parallel()
	testsize.Small(t)

	theme := New(colorprofile.TrueColor)
	p := theme.Colors
	tests := []struct {
		name     string
		status   model.TaskStatus
		selected bool
		want     color.Color
	}{
		{"ready selected", model.TaskStatusReady, true, p.FG},
		{"ready unselected", model.TaskStatusReady, false, p.FG},
		{"running selected", model.TaskStatusRunning, true, p.Success},
		{"running unselected", model.TaskStatusRunning, false, p.SuccessMuted},
		{"pending feedback selected", model.TaskStatusPendingFeedback, true, p.Warning},
		{"pending feedback unselected", model.TaskStatusPendingFeedback, false, p.WarningMuted},
		{"completed selected", model.TaskStatusCompleted, true, p.FGMuted},
		{"completed unselected", model.TaskStatusCompleted, false, p.FGMuted},
		{"canceled selected", model.TaskStatusCanceled, true, p.FGMuted},
		{"canceled unselected", model.TaskStatusCanceled, false, p.FGMuted},
		{"failed selected", model.TaskStatusFailed, true, p.Danger},
		{"failed unselected", model.TaskStatusFailed, false, p.DangerMuted},
		{"unknown selected", model.TaskStatus("unknown"), true, p.FG},
		{"unknown unselected", model.TaskStatus("unknown"), false, p.FG},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			label := theme.StatusLabel(tc.selected, tc.status)
			if label.GetForeground() != tc.want {
				t.Fatalf("status label foreground = %v, want %v", label.GetForeground(), tc.want)
			}
			assertNoBackground(t, label, "status label")
			if label.GetPaddingLeft() != 0 || label.GetPaddingRight() != 0 {
				t.Fatalf("status label padding = (%d, %d), want (0, 0)", label.GetPaddingLeft(), label.GetPaddingRight())
			}
			if label.GetBorderLeft() {
				t.Fatal("status label has a left border, want none")
			}
		})
	}
}

func TestTaskItemsLayoutIndependentOfStatus(t *testing.T) {
	t.Parallel()
	testsize.Small(t)

	theme := New(colorprofile.TrueColor)
	base := theme.TaskItems(model.TaskStatusReady)
	statuses := []model.TaskStatus{
		model.TaskStatusReady,
		model.TaskStatusRunning,
		model.TaskStatusPendingFeedback,
		model.TaskStatusCompleted,
		model.TaskStatusCanceled,
		model.TaskStatusFailed,
	}
	for _, status := range statuses {
		t.Run(status.String(), func(t *testing.T) {
			t.Parallel()
			items := theme.TaskItems(status)
			if items.NormalTitle.GetPaddingLeft() != base.NormalTitle.GetPaddingLeft() {
				t.Fatalf("normal title padding left = %d, want %d", items.NormalTitle.GetPaddingLeft(), base.NormalTitle.GetPaddingLeft())
			}
			if items.SelectedTitle.GetPaddingLeft() != base.SelectedTitle.GetPaddingLeft() {
				t.Fatalf("selected title padding left = %d, want %d", items.SelectedTitle.GetPaddingLeft(), base.SelectedTitle.GetPaddingLeft())
			}
			if items.SelectedDesc.GetPaddingLeft() != base.SelectedDesc.GetPaddingLeft() {
				t.Fatalf("selected description padding left = %d, want %d", items.SelectedDesc.GetPaddingLeft(), base.SelectedDesc.GetPaddingLeft())
			}
			if items.SelectedTitle.GetBorderStyle().Left != selectedMarker {
				t.Fatalf("selected title glyph = %q, want %q", items.SelectedTitle.GetBorderStyle().Left, selectedMarker)
			}
			if items.SelectedTitle.GetBorderLeftForeground() != theme.Colors.FG {
				t.Fatalf("selected title left border = %v, want fg %v", items.SelectedTitle.GetBorderLeftForeground(), theme.Colors.FG)
			}
			if items.SelectedDesc.GetBorderLeft() {
				t.Fatal("selected description has a left border, want padding only")
			}
		})
	}
}

func assertNoBackground(t *testing.T, s lipgloss.Style, name string) {
	t.Helper()
	if _, ok := s.GetBackground().(lipgloss.NoColor); !ok {
		t.Fatalf("%s background = %v, want no color", name, s.GetBackground())
	}
}

func hasBoldSGR(s string) bool {
	return strings.Contains(s, "\x1b[1m") || strings.Contains(s, "\x1b[1;") || strings.Contains(s, ";1;") || strings.Contains(s, ";1m")
}

func assertSelectedTitleMarker(t *testing.T, s lipgloss.Style, p Palette, name string) {
	t.Helper()
	if !s.GetBorderLeft() || s.GetBorderLeftSize() == 0 {
		t.Fatalf("%s has no left border", name)
	}
	if s.GetBorderTop() || s.GetBorderRight() || s.GetBorderBottom() {
		t.Fatalf("%s has extra borders, want left only", name)
	}
	if s.GetBorderStyle().Left != selectedMarker {
		t.Fatalf("%s left glyph = %q, want %q", name, s.GetBorderStyle().Left, selectedMarker)
	}
	if s.GetBorderLeftForeground() != p.FG {
		t.Fatalf("%s left border = %v, want fg %v", name, s.GetBorderLeftForeground(), p.FG)
	}
	if s.GetBorderLeftSize()+s.GetPaddingLeft() != commonPad {
		t.Fatalf("%s left inset = %d, want %d", name, s.GetBorderLeftSize()+s.GetPaddingLeft(), commonPad)
	}
}

func assertPaddedNoBorder(t *testing.T, s lipgloss.Style, name string) {
	t.Helper()
	if s.GetBorderLeft() || s.GetBorderLeftSize() != 0 {
		t.Fatalf("%s has a left border, want padding only", name)
	}
	if s.GetPaddingLeft() != commonPad {
		t.Fatalf("%s padding left = %d, want %d", name, s.GetPaddingLeft(), commonPad)
	}
}
