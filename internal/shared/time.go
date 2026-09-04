package shared

import (
	"fmt"
	"time"
)

// RelativeTime returns a compact duration like "2s ago" or "1d ago".
// A zero time returns an empty string. Future times return "0s ago".
func RelativeTime(t, now time.Time) string {
	if t.IsZero() {
		return ""
	}
	d := max(now.Sub(t), 0)
	switch {
	case d < time.Minute:
		return fmt.Sprintf("%ds ago", int(d/time.Second))
	case d < time.Hour:
		return fmt.Sprintf("%dm ago", int(d/time.Minute))
	case d < 24*time.Hour:
		return fmt.Sprintf("%dh ago", int(d/time.Hour))
	default:
		return fmt.Sprintf("%dd ago", int(d/(24*time.Hour)))
	}
}
