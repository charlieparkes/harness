package styles

import (
	"image/color"
	"testing"

	"charm.land/lipgloss/v2"
	"github.com/charlieparkes/go-testsize"
	"github.com/charmbracelet/colorprofile"
)

func TestTokenColorsMatchCharmDefaults(t *testing.T) {
	t.Parallel()
	testsize.Small(t)

	want := map[ColorType]color.Color{
		ColorTypeBg:           lipgloss.Color("#3C3C3C"),
		ColorTypeFg:           lipgloss.Color("#DDDDDD"),
		ColorTypeFgMuted:      lipgloss.Color("#777777"),
		ColorTypeAccent:       lipgloss.Color("#EE6FF8"),
		ColorTypeBorder:       lipgloss.Color("#5C5C5C"),
		ColorTypeSuccess:      lipgloss.Color("#87FF00"),
		ColorTypeSuccessMuted: lipgloss.Color("#87D75F"),
		ColorTypeWarning:      lipgloss.Color("#FF8700"),
		ColorTypeWarningMuted: lipgloss.Color("#D7AF87"),
		ColorTypeDanger:       lipgloss.Color("#FF005F"),
		ColorTypeDangerMuted:  lipgloss.Color("#FF5F87"),
	}
	defs := tokens()
	if len(defs) != len(want) {
		t.Fatalf("len(tokens) = %d, want %d", len(defs), len(want))
	}
	for _, tok := range defs {
		c, ok := want[tok.name]
		if !ok {
			t.Fatalf("unexpected token %q", tok.name)
		}
		if tok.trueColor != c {
			t.Fatalf("%s = %v (%T), want %v (%T)", tok.name, tok.trueColor, tok.trueColor, c, c)
		}
		delete(want, tok.name)
	}
	if len(want) != 0 {
		t.Fatalf("missing tokens %v", want)
	}
}

func TestTokensResolveForProfiles(t *testing.T) {
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
			p := resolve(profile)
			for _, tok := range tokens() {
				got, ok := p.color(tok.name)
				if !ok {
					t.Fatalf("missing token %q", tok.name)
				}
				want := wantColor(profile, tok)
				if got != want {
					t.Fatalf("%s = %v (%T), want %v (%T)", tok.name, got, got, want, want)
				}
			}
		})
	}
}

func TestANSIFallbacks(t *testing.T) {
	t.Parallel()
	testsize.Small(t)

	want := map[ColorType]color.Color{
		ColorTypeBg:           lipgloss.Black,
		ColorTypeFg:           lipgloss.White,
		ColorTypeFgMuted:      lipgloss.BrightBlack,
		ColorTypeAccent:       lipgloss.BrightMagenta,
		ColorTypeBorder:       lipgloss.BrightBlack,
		ColorTypeSuccess:      lipgloss.BrightGreen,
		ColorTypeSuccessMuted: lipgloss.Green,
		ColorTypeWarning:      lipgloss.BrightYellow,
		ColorTypeWarningMuted: lipgloss.Yellow,
		ColorTypeDanger:       lipgloss.BrightRed,
		ColorTypeDangerMuted:  lipgloss.Red,
	}
	p := resolve(colorprofile.ANSI)
	defs := tokens()
	if len(defs) != len(want) {
		t.Fatalf("len(tokens) = %d, want %d", len(defs), len(want))
	}
	for _, tok := range defs {
		ansi, ok := want[tok.name]
		if !ok {
			t.Fatalf("unexpected token %q", tok.name)
		}
		if tok.ansi != ansi {
			t.Fatalf("%s ansi = %v (%T), want %v (%T)", tok.name, tok.ansi, tok.ansi, ansi, ansi)
		}
		got, ok := p.color(tok.name)
		if !ok {
			t.Fatalf("missing token %q", tok.name)
		}
		if got != ansi {
			t.Fatalf("%s = %v (%T), want %v (%T)", tok.name, got, got, ansi, ansi)
		}
	}
}

func TestANSI256Fallbacks(t *testing.T) {
	t.Parallel()
	testsize.Small(t)

	want := map[ColorType]color.Color{
		ColorTypeBg:           lipgloss.ANSIColor(234),
		ColorTypeFg:           lipgloss.ANSIColor(253),
		ColorTypeFgMuted:      lipgloss.ANSIColor(243),
		ColorTypeAccent:       lipgloss.ANSIColor(207),
		ColorTypeBorder:       lipgloss.ANSIColor(59),
		ColorTypeSuccess:      lipgloss.ANSIColor(118),
		ColorTypeSuccessMuted: lipgloss.ANSIColor(113),
		ColorTypeWarning:      lipgloss.ANSIColor(208),
		ColorTypeWarningMuted: lipgloss.ANSIColor(180),
		ColorTypeDanger:       lipgloss.ANSIColor(197),
		ColorTypeDangerMuted:  lipgloss.ANSIColor(204),
	}
	p := resolve(colorprofile.ANSI256)
	defs := tokens()
	if len(defs) != len(want) {
		t.Fatalf("len(tokens) = %d, want %d", len(defs), len(want))
	}
	for _, tok := range defs {
		ansi256, ok := want[tok.name]
		if !ok {
			t.Fatalf("unexpected token %q", tok.name)
		}
		if tok.ansi256 != ansi256 {
			t.Fatalf("%s ansi256 = %v (%T), want %v (%T)", tok.name, tok.ansi256, tok.ansi256, ansi256, ansi256)
		}
		got, ok := p.color(tok.name)
		if !ok {
			t.Fatalf("missing token %q", tok.name)
		}
		if got != ansi256 {
			t.Fatalf("%s = %v (%T), want %v (%T)", tok.name, got, got, ansi256, ansi256)
		}
	}
}

func TestDangerANSI256IsNotGreyscale(t *testing.T) {
	t.Parallel()
	testsize.Small(t)

	p := resolve(colorprofile.ANSI256)
	idx, ok := p.Danger.(lipgloss.ANSIColor)
	if !ok {
		t.Fatalf("danger type = %T, want ANSIColor", p.Danger)
	}
	if idx >= 232 {
		t.Fatalf("danger ansi256 = %d (greyscale), want a red cube color", idx)
	}
	muted, ok := p.FGMuted.(lipgloss.ANSIColor)
	if !ok {
		t.Fatalf("fg-muted type = %T, want ANSIColor", p.FGMuted)
	}
	if idx == muted {
		t.Fatalf("danger ansi256 = %d, same as fg-muted", idx)
	}
}

func wantColor(profile colorprofile.Profile, tok tokenDef) color.Color {
	switch profile {
	case colorprofile.TrueColor:
		return tok.trueColor
	case colorprofile.ANSI256:
		return tok.ansi256
	case colorprofile.ANSI:
		return tok.ansi
	case colorprofile.Unknown, colorprofile.NoTTY, colorprofile.ASCII:
		return lipgloss.NoColor{}
	default:
		return lipgloss.NoColor{}
	}
}
