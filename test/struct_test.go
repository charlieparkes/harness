package test

import (
	"testing"

	"github.com/charlieparkes/go-testsize"
)

func TestAssertStructFieldsEqual(t *testing.T) {
	t.Parallel()
	testsize.Small(t)

	type got struct {
		ID    int
		Title string
	}
	type want struct {
		ID    string
		Title bool
	}

	AssertStructFieldsEqual(t, got{}, want{})
	AssertStructFieldsEqual(t, &got{}, want{})
}

func TestAssertStructFieldsEqualIgnore(t *testing.T) {
	t.Parallel()
	testsize.Small(t)

	type got struct {
		ID     int
		Title  string
		Timing string
	}
	type want struct {
		ID    string
		Title bool
		Extra bool
	}

	AssertStructFieldsEqual(t, got{}, want{}, "Timing", "Extra")
}
