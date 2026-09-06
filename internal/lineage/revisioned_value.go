package lineage

type RevisionedValue[T any] struct {
	Revision int64 `json:"revision"` // parent document revision number at which this value was set
	Value    T     `json:"value,omitempty"`
}
