package lineage

import (
	"slices"
)

// RevisionedField represents all values of type FieldType for field FieldName.
type RevisionedField[FieldType any] struct {
	Values []RevisionedValue[FieldType] `json:"values,omitempty"`
}

func NewRevisionedField[FieldType any](values ...FieldType) RevisionedField[FieldType] {
	var revisionedValues []RevisionedValue[FieldType]
	if len(values) > 0 {
		revisionedValues = make([]RevisionedValue[FieldType], len(values))
		for i, v := range values {
			revisionedValues[int64(i)] = NewRevisionedValue(int64(i), v)
		}
	}
	return RevisionedField[FieldType]{
		Values: revisionedValues,
	}
}

// Revisions returns an array of revision numbers at which [RevisionedField] stores a value.
func (f RevisionedField[FieldType]) Revisions() []int64 {
	var revisions []int64
	for _, v := range f.Values {
		revisions = append(revisions, v.Revision)
	}
	return revisions
}

// Value returns the latest value of a [RevisionedField].
func (f RevisionedField[FieldType]) Value() (FieldType, bool) {
	if len(f.Values) == 0 {
		var zero FieldType
		return zero, false
	}
	return f.Values[slices.Max(f.Revisions())].Value, true
}

// ValueAt returns the most recent value of [RevisionedField] in relation to a given revision.
// If requested revision is older than the first known value, return false.
func (f RevisionedField[FieldType]) ValueAt(revision int64) (FieldType, bool) {
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

// Set value at particular revision.
// Must be newer than any known revision.
// Returns bool indicating success.
func (f RevisionedField[FieldType]) Set(revision int64, value FieldType) bool {
	if revision <= slices.Max(f.Revisions()) {
		return false
	}
	f.Values = append(f.Values, NewRevisionedValue(revision, value))
	return true
}
