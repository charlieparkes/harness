package model

import (
	"time"

	"github.com/charlieparkes/harness/internal/lineage"
)

// Plan is the plan document for a task.
type Plan struct {
	ID       string `json:"id"`
	Revision int64  `json:"revision"`
	TaskID   string `json:"task_id"`

	Title            lineage.RevisionedField[string]     `json:"title"`
	Description      lineage.RevisionedField[string]     `json:"description"`
	Status           lineage.RevisionedField[PlanStatus] `json:"status"`
	Risks            lineage.RevisionedField[[]string]   `json:"risks"`
	DefinitionOfDone lineage.RevisionedField[[]string]   `json:"definition_of_done"`
	Steps            lineage.RevisionedField[[]Step]     `json:"steps"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}
