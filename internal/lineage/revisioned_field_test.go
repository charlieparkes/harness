package lineage

import (
	"slices"
	"testing"

	"github.com/charlieparkes/go-testcmp"
	"github.com/charlieparkes/go-testsize"
	"github.com/stretchr/testify/require"
)

func TestNewRevisionedField(t *testing.T) {
	t.Parallel()
	testsize.Small(t)

	t.Run("no values", func(t *testing.T) {
		t.Parallel()
		testcmp.Compare(t, NewRevisionedField[int](), RevisionedField[int]{})
	})

	t.Run("records values at sequential revisions", func(t *testing.T) {
		t.Parallel()
		testcmp.Compare(t, NewRevisionedField(10, 20, 30), RevisionedField[int]{
			Values: []RevisionedValue[int]{
				{Revision: 1, Value: 10},
				{Revision: 2, Value: 20},
				{Revision: 3, Value: 30},
			},
		})
	})

	t.Run("skips leading zero value", func(t *testing.T) {
		t.Parallel()
		testcmp.Compare(t, NewRevisionedField(0, 10), RevisionedField[int]{
			Values: []RevisionedValue[int]{
				{Revision: 1, Value: 10},
			},
		})
	})

	t.Run("skips nil slice", func(t *testing.T) {
		t.Parallel()
		testcmp.Compare(t, NewRevisionedField[[]string](nil), RevisionedField[[]string]{})
	})

	t.Run("skips equal consecutive values", func(t *testing.T) {
		t.Parallel()
		testcmp.Compare(t, NewRevisionedField(10, 10, 20), RevisionedField[int]{
			Values: []RevisionedValue[int]{
				{Revision: 1, Value: 10},
				{Revision: 2, Value: 20},
			},
		})
	})

	t.Run("records a later zero when it differs from the previous value", func(t *testing.T) {
		t.Parallel()
		testcmp.Compare(t, NewRevisionedField(10, 0, 20), RevisionedField[int]{
			Values: []RevisionedValue[int]{
				{Revision: 1, Value: 10},
				{Revision: 2, Value: 0},
				{Revision: 3, Value: 20},
			},
		})
	})
}

func TestRevisionedFieldValue(t *testing.T) {
	t.Parallel()
	testsize.Small(t)

	tests := []struct {
		name  string
		field RevisionedField[int]
		want  int
	}{
		{
			name:  "empty",
			field: RevisionedField[int]{},
		},
		{
			name: "single value",
			field: RevisionedField[int]{
				Values: []RevisionedValue[int]{
					{Revision: 2, Value: 20},
				},
			},
			want: 20,
		},
		{
			name:  "latest of several",
			field: revisionedIntField(),
			want:  50,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, tt.field.Value())
		})
	}
}

func TestRevisionedFieldValueAt(t *testing.T) {
	t.Parallel()
	testsize.Small(t)

	field := revisionedIntField()
	empty := RevisionedField[int]{}

	tests := []struct {
		name     string
		field    RevisionedField[int]
		revision int64
		want     int
	}{
		{name: "empty field", field: empty, revision: 1},
		{name: "zero revision", field: field, revision: 0},
		{name: "negative revision", field: field, revision: -1},
		{name: "before first value", field: field, revision: 1},
		{name: "at first value", field: field, revision: 2, want: 20},
		{name: "between values", field: field, revision: 3, want: 20},
		{name: "at middle value", field: field, revision: 4, want: 40},
		{name: "after middle value", field: field, revision: 5, want: 40},
		{name: "at latest value", field: field, revision: 7, want: 50},
		{name: "after latest value", field: field, revision: 8, want: 50},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, tt.field.ValueAt(tt.revision))
		})
	}
}

func TestRevisionedFieldSet(t *testing.T) {
	t.Parallel()
	testsize.Small(t)

	t.Run("records on empty field", func(t *testing.T) {
		t.Parallel()
		var field RevisionedField[int]
		if ok := field.Set(1, 10); !ok {
			t.Fatalf("Set(1, 10) = false, want true")
		}
		require.Equal(t, 10, field.Value())
		if got := field.Revisions(); !slices.Equal(got, []int64{1}) {
			t.Fatalf("Revisions() = %v, want [1]", got)
		}
	})

	t.Run("skips zero value on empty field", func(t *testing.T) {
		t.Parallel()
		var field RevisionedField[[]string]
		if ok := field.Set(1, nil); ok {
			t.Fatalf("Set(1, nil) = true, want false")
		}
		if got := field.Revisions(); len(got) != 0 {
			t.Fatalf("Revisions() = %v, want []", got)
		}
	})

	t.Run("records newer revision with changed value", func(t *testing.T) {
		t.Parallel()
		field := RevisionedField[int]{
			Values: []RevisionedValue[int]{{Revision: 1, Value: 10}},
		}
		if ok := field.Set(2, 20); !ok {
			t.Fatalf("Set(2, 20) = false, want true")
		}
		require.Equal(t, 20, field.Value())
		if got := field.Revisions(); !slices.Equal(got, []int64{1, 2}) {
			t.Fatalf("Revisions() = %v, want [1 2]", got)
		}
	})

	t.Run("skips newer revision with equal value", func(t *testing.T) {
		t.Parallel()
		field := RevisionedField[[]string]{
			Values: []RevisionedValue[[]string]{{Revision: 1, Value: []string{"a", "b"}}},
		}
		if ok := field.Set(2, []string{"a", "b"}); ok {
			t.Fatalf("Set(2, [a b]) = true, want false")
		}
		if got := field.Revisions(); !slices.Equal(got, []int64{1}) {
			t.Fatalf("Revisions() = %v, want [1]", got)
		}
	})

	t.Run("rejects revision equal to or older than latest", func(t *testing.T) {
		t.Parallel()
		field := RevisionedField[int]{
			Values: []RevisionedValue[int]{{Revision: 2, Value: 20}},
		}
		if ok := field.Set(2, 30); ok {
			t.Fatalf("Set(2, 30) = true, want false")
		}
		if ok := field.Set(1, 30); ok {
			t.Fatalf("Set(1, 30) = true, want false")
		}
		if got := field.Revisions(); !slices.Equal(got, []int64{2}) {
			t.Fatalf("Revisions() = %v, want [2]", got)
		}
	})
}

func revisionedIntField() RevisionedField[int] {
	return RevisionedField[int]{
		Values: []RevisionedValue[int]{
			{Revision: 2, Value: 20},
			{Revision: 4, Value: 40},
			{Revision: 7, Value: 50},
		},
	}
}
