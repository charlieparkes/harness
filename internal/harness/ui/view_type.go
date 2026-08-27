package ui

//go:generate go tool go-enum --nocamel

// ENUM(tasks, example)
type ViewType int

// Name returns the display name for the view type.
func (v ViewType) Name() string {
	switch v {
	case ViewTypeTasks:
		return tasksTitle
	case ViewTypeExample:
		return exampleTitle
	default:
		return ""
	}
}
