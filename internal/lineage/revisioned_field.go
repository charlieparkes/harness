package lineage

import (
	"reflect"
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
			revisionedValues[i] = NewRevisionedValue(int64(i)+1, v)
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
	return f.Values[len(f.Values)-1].Value, true
}

// ValueAt returns the most recent value of [RevisionedField] in relation to a given revision.
// If requested revision is older than the first known value, return false.
func (f RevisionedField[FieldType]) ValueAt(revision int64) (FieldType, bool) {
	if revision <= 0 || len(f.Values) == 0 {
		var zero FieldType
		return zero, false
	}
	if revision < f.Values[0].Revision {
		var zero FieldType
		return zero, false
	}
	value := f.Values[0].Value
	for _, v := range f.Values {
		if v.Revision > revision {
			break
		}
		value = v.Value
	}
	return value, true
}

// Set records value at a particular revision.
// The revision must be newer than any known revision, and the value must
// differ from the latest value. Returns true only when a value was recorded.
func (f *RevisionedField[FieldType]) Set(revision int64, value FieldType) bool {
	if len(f.Values) > 0 {
		last := f.Values[len(f.Values)-1]
		if revision <= last.Revision {
			return false
		}
		if reflect.DeepEqual(last.Value, value) {
			return false
		}
	}
	f.Values = append(f.Values, NewRevisionedValue(revision, value))
	return true
}
