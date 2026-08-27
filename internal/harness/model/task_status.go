package model

//go:generate go tool go-enum

// TaskStatus is the lifecycle state of a task.
// ENUM(Unspecified, Ready, Running, Pending Feedback, Completed, Canceled, Failed)
type TaskStatus string

// Symbol returns the display glyph for this status.
func (s TaskStatus) Symbol() string {
	switch s {
	case TaskStatusUnspecified:
		return ""
	case TaskStatusReady:
		return "○"
	case TaskStatusRunning:
		return "●"
	case TaskStatusPendingFeedback:
		return "◐"
	case TaskStatusCompleted:
		return "✓"
	case TaskStatusCanceled:
		return "⦸"
	case TaskStatusFailed:
		return "✕"
	default:
		return ""
	}
}

// Text returns the status label: glyph then English name.
func (s TaskStatus) Text() string {
	name := s.String()
	if name == "" {
		return ""
	}
	if sym := s.Symbol(); sym != "" {
		return sym + " " + name
	}
	return name
}
