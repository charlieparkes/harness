package model

type File struct {
	TargetID string `json:"target_id"`
	Path     string `json:"path"`
	Line     *int64 `json:"line"`
	Exists   bool   `json:"exists"`
}
