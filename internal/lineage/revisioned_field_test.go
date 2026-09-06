package lineage

import (
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
					2: {Revision: 2, Value: 20},
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
		{name: "before first value", field: field, revision: 1, want: 0, ok: true},
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

func revisionedIntField() RevisionedField[int] {
	return RevisionedField[int]{
		Values: []RevisionedValue[int]{
			2: {Revision: 2, Value: 20},
			4: {Revision: 4, Value: 40},
			7: {Revision: 7, Value: 50},
		},
	}
}
