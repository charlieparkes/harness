package lineage

type RevisionedValue[T any] struct {
	Revision int64 `json:"revision"` // parent document revision number at which this value was set
	Value    T     `json:"value,omitempty"`
}

func NewRevisionedValue[T any](revision int64, value T) RevisionedValue[T] {
	return RevisionedValue[T]{
		Revision: revision,
		Value:    value,
	}
}
