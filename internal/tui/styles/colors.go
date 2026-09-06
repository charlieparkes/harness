package styles

import (
	"image/color"
	"sync"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/colorprofile"
)

// Palette is the TUI colors resolved for a terminal profile.
type Palette struct {
	BG           color.Color
	FG           color.Color
	FGMuted      color.Color
	Accent       color.Color
	Border       color.Color
	Success      color.Color
	SuccessMuted color.Color
	Warning      color.Color
	WarningMuted color.Color
	Danger       color.Color
	DangerMuted  color.Color
}

type tokenDef struct {
	name      ColorType
	trueColor color.Color
	ansi256   color.Color
	ansi      color.Color
}

var tokens = sync.OnceValue(func() []tokenDef {
	// Note: Comment containing hex code for each trueColor, ansi256, and ansi
	// so that supporting IDEs will display the color picker.
	return []tokenDef{
		{
			name:      ColorTypeBg,
			trueColor: lipgloss.Color("#3C3C3C"), // #3C3C3C
			ansi256:   lipgloss.ANSIColor(234),   // #3A3A3A
			ansi:      lipgloss.Black,            // #000000
		},
		{
			name:      ColorTypeFg,
			trueColor: lipgloss.Color("#DDDDDD"), // #DDDDDD
			ansi256:   lipgloss.ANSIColor(253),   // #DADADA
			ansi:      lipgloss.White,            // #C0C0C0
		},
		{
			name:      ColorTypeFgMuted,
			trueColor: lipgloss.Color("#777777"), // #777777
			ansi256:   lipgloss.ANSIColor(243),   // #767676
			ansi:      lipgloss.BrightBlack,      // #808080
		},
		{
			name:      ColorTypeAccent,
			trueColor: lipgloss.Color("#EE6FF8"), // #EE6FF8
			ansi256:   lipgloss.ANSIColor(207),   // #FF5FFF
			ansi:      lipgloss.BrightMagenta,    // #FF00FF
		},
		{
			name:      ColorTypeBorder,
			trueColor: lipgloss.Color("#5C5C5C"), // #5C5C5C
			ansi256:   lipgloss.ANSIColor(59),    // #5F5F5F
			ansi:      lipgloss.BrightBlack,      // #808080
		},
		{
			name:      ColorTypeSuccess,
			trueColor: lipgloss.Color("#87FF00"), // #87FF00
			ansi256:   lipgloss.ANSIColor(118),   // #87FF00
			ansi:      lipgloss.BrightGreen,      // #00FF00
		},
		{
			name:      ColorTypeSuccessMuted,
			trueColor: lipgloss.Color("#87D75F"), // #87D75F
			ansi256:   lipgloss.ANSIColor(113),   // #87D75F
			ansi:      lipgloss.Green,            // #008000
		},
		{
			name:      ColorTypeWarning,
			trueColor: lipgloss.Color("#FF8700"), // #FF8700
			ansi256:   lipgloss.ANSIColor(208),   // #FF8700
			ansi:      lipgloss.BrightYellow,     // #FFFF00
		},
		{
			name:      ColorTypeWarningMuted,
			trueColor: lipgloss.Color("#D7AF87"), // #D7AF87
			ansi256:   lipgloss.ANSIColor(180),   // #D7AF87
			ansi:      lipgloss.Yellow,           // #808000
		},
		{
			name:      ColorTypeDanger,
			trueColor: lipgloss.Color("#FF005F"), // #FF005F
			ansi256:   lipgloss.ANSIColor(197),   // #FF005F
			ansi:      lipgloss.BrightRed,        // #FF0000
		},
		{
			name:      ColorTypeDangerMuted,
			trueColor: lipgloss.Color("#FF5F87"), // #FF5F87
			ansi256:   lipgloss.ANSIColor(204),   // #FF5F87
			ansi:      lipgloss.Red,              // #800000
		},
	}
})

func resolve(profile colorprofile.Profile) Palette {
	complete := lipgloss.Complete(profile)
	defs := tokens()
	colors := make(map[ColorType]color.Color, len(defs))
	for _, tok := range defs {
		colors[tok.name] = complete(tok.ansi, tok.ansi256, tok.trueColor)
	}
	return Palette{
		BG:           colors[ColorTypeBg],
		FG:           colors[ColorTypeFg],
		FGMuted:      colors[ColorTypeFgMuted],
		Accent:       colors[ColorTypeAccent],
		Border:       colors[ColorTypeBorder],
		Success:      colors[ColorTypeSuccess],
		SuccessMuted: colors[ColorTypeSuccessMuted],
		Warning:      colors[ColorTypeWarning],
		WarningMuted: colors[ColorTypeWarningMuted],
		Danger:       colors[ColorTypeDanger],
		DangerMuted:  colors[ColorTypeDangerMuted],
	}
}

func (p Palette) color(name ColorType) (color.Color, bool) {
	colors := map[ColorType]color.Color{
		ColorTypeBg:           p.BG,
		ColorTypeFg:           p.FG,
		ColorTypeFgMuted:      p.FGMuted,
		ColorTypeAccent:       p.Accent,
		ColorTypeBorder:       p.Border,
		ColorTypeSuccess:      p.Success,
		ColorTypeSuccessMuted: p.SuccessMuted,
		ColorTypeWarning:      p.Warning,
		ColorTypeWarningMuted: p.WarningMuted,
		ColorTypeDanger:       p.Danger,
		ColorTypeDangerMuted:  p.DangerMuted,
	}
	c, ok := colors[name]
	return c, ok
}
