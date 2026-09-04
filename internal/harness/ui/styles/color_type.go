package styles

//go:generate go tool go-enum

// ColorType is a palette token.
// ENUM(bg, fg, fg-muted, accent, border, success, success-muted, warning, warning-muted, danger, danger-muted)
type ColorType string
