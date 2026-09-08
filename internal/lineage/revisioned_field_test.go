package lineage

import (
	"slices"
	"testing"

	"github.com/charlieparkes/go-testsize"
)

func TestRevisionedFieldValue(t *testing.T) {
	t.Parallel()
	testsize.Small(t)

	tests := []struct {
		name  string
		field RevisionedField[int]
		want  int
		ok    bool
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
			ok:   true,
		},
		{
			name:  "latest of several",
			field: revisionedIntField(),
			want:  50,
			ok:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, ok := tt.field.Value()
			if ok != tt.ok || got != tt.want {
				t.Fatalf("Value() = (%d, %t), want (%d, %t)", got, ok, tt.want, tt.ok)
			}
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
		ok       bool
	}{
		{name: "empty field", field: empty, revision: 1},
		{name: "zero revision", field: field, revision: 0},
		{name: "negative revision", field: field, revision: -1},
		{name: "before first value", field: field, revision: 1},
		{name: "at first value", field: field, revision: 2, want: 20, ok: true},
		{name: "between values", field: field, revision: 3, want: 20, ok: true},
		{name: "at middle value", field: field, revision: 4, want: 40, ok: true},
		{name: "after middle value", field: field, revision: 5, want: 40, ok: true},
		{name: "at latest value", field: field, revision: 7, want: 50, ok: true},
		{name: "after latest value", field: field, revision: 8, want: 50, ok: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, ok := tt.field.ValueAt(tt.revision)
			if ok != tt.ok || got != tt.want {
				t.Fatalf("ValueAt(%d) = (%d, %t), want (%d, %t)", tt.revision, got, ok, tt.want, tt.ok)
			}
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
		if got, ok := field.Value(); got != 10 || !ok {
			t.Fatalf("Value() = (%d, %t), want (10, true)", got, ok)
		}
		if got := field.Revisions(); !slices.Equal(got, []int64{1}) {
			t.Fatalf("Revisions() = %v, want [1]", got)
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
		if got, _ := field.Value(); got != 20 {
			t.Fatalf("Value() = %d, want 20", got)
		}
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
