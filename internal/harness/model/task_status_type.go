package model

import (
	"testing"

	"github.com/charlieparkes/go-testsize"
)

func TestStatusString(t *testing.T) {
	t.Parallel()
	testsize.Small(t)

	tests := []struct {
		status TaskStatus
		want   string
	}{
		{TaskStatusReady, "Ready"},
		{TaskStatusRunning, "Running"},
		{TaskStatusPendingFeedback, "Pending Feedback"},
		{TaskStatusCompleted, "Completed"},
		{TaskStatusCanceled, "Canceled"},
		{TaskStatusFailed, "Failed"},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			t.Parallel()
			if got := tt.status.String(); got != tt.want {
				t.Fatalf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestStatusText(t *testing.T) {
	t.Parallel()
	testsize.Small(t)

	tests := []struct {
		name   string
		status TaskStatus
		symbol string
		text   string
	}{
		{name: "ready", status: TaskStatusReady, symbol: "○", text: "○ Ready"},
		{name: "running", status: TaskStatusRunning, symbol: "●", text: "● Running"},
		{name: "pending feedback", status: TaskStatusPendingFeedback, symbol: "◐", text: "◐ Pending Feedback"},
		{name: "completed", status: TaskStatusCompleted, symbol: "✓", text: "✓ Completed"},
		{name: "canceled", status: TaskStatusCanceled, symbol: "⦸", text: "⦸ Canceled"},
		{name: "failed", status: TaskStatusFailed, symbol: "✕", text: "✕ Failed"},
		{name: "empty", status: "", symbol: "", text: ""},
		{name: "unknown", status: TaskStatus("unknown"), symbol: "", text: "unknown"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := tc.status.Symbol(); got != tc.symbol {
				t.Fatalf("Symbol() = %q, want %q", got, tc.symbol)
			}
			if got := tc.status.Text(); got != tc.text {
				t.Fatalf("Text() = %q, want %q", got, tc.text)
			}
		})
	}
}
