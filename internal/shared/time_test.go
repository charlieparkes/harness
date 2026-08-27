package shared

import (
	"testing"
	"time"

	"github.com/charlieparkes/go-testsize"
)

func TestRelativeTime(t *testing.T) {
	t.Parallel()
	testsize.Small(t)

	now := time.Date(2026, 8, 30, 12, 0, 0, 0, time.UTC)
	tests := []struct {
		name string
		at   time.Time
		want string
	}{
		{name: "zero", want: ""},
		{name: "same instant", at: now, want: "0s ago"},
		{name: "future", at: now.Add(time.Hour), want: "0s ago"},
		{name: "seconds", at: now.Add(-2 * time.Second), want: "2s ago"},
		{name: "minutes", at: now.Add(-3 * time.Minute), want: "3m ago"},
		{name: "hours", at: now.Add(-5 * time.Hour), want: "5h ago"},
		{name: "days", at: now.Add(-26 * time.Hour), want: "1d ago"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := RelativeTime(tt.at, now)
			if got != tt.want {
				t.Fatalf("RelativeTime() = %q, want %q", got, tt.want)
			}
		})
	}
}
