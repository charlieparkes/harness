package lineage

import (
	"slices"
)

// RevisionedField represents all values of type FieldType for field FieldName.
type RevisionedField[FieldName comparable, FieldType any] struct {
	Name   FieldName                            `json:"name"` // enum of supported field names
	Values map[int64]RevisionedValue[FieldType] `json:"values,omitempty"`
}

// Revisions returns an array of revision numbers at which [RevisionedField] stores a value.
func (f RevisionedField[FieldName, FieldType]) Revisions() []int64 {
	var revisions []int64
	for k := range f.Values {
		revisions = append(revisions, k)
	}
	return revisions
}

// Value returns the latest value of a [RevisionedField].
func (f RevisionedField[FieldName, FieldType]) Value() (FieldType, bool) {
	if len(f.Values) == 0 {
		var zero FieldType
		return zero, false
	}
	return f.Values[slices.Max(f.Revisions())].Value, true
}

// ValueAt returns the most recent value of [RevisionedField] in relation to a given revision.
// If requested revision is older than the first known value, return false.
func (f RevisionedField[FieldName, FieldType]) ValueAt(revision int64) (FieldType, bool) {
	if revision <= 0 || len(f.Values) == 0 {
		var zero FieldType
		return zero, false
	}
	revisions := f.Revisions()
	oldest := slices.Min(revisions)
	if oldest > revision {
		var zero FieldType
		return zero, false
	}
	latest := slices.Max(revisions)
	if latest > revision {
		revisions = slices.DeleteFunc(revisions, func(r int64) bool {
			return r > revision
		})
		latest = slices.Max(revisions)
	}
	return f.Values[latest].Value, true
}
