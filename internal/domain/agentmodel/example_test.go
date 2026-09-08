package agentmodel

import (
	"testing"
	"time"

	"github.com/charlieparkes/go-testcmp"
	"github.com/charlieparkes/go-testsize"
	"github.com/charlieparkes/harness/internal/domain/model"
	"github.com/charlieparkes/harness/internal/lineage"
)

// This is a human-written example of the expected behavior when translating between
// domain models and agent models.
func TestExampleUsage(t *testing.T) {
	t.Parallel()
	_ = testsize.Small(t)
	now := time.Now().UTC().Truncate(time.Second)
	planID := "PLAN-0001"

	t.Run("Agent Creates First Revision", func(t *testing.T) {
		t.Parallel()
		newPlan := Plan{
			Title:       "Plan Title",
			Description: "Plan Description",
			Steps: []Step{
				{
					ID:      "S1",
					Title:   "Step Title",
					Summary: "Step Summary",
					ProposedChanges: []ProposedChange{
						{
							Description: "Proposed Change Description",
							Reason:      "Proposed Change Reason",
							Files: []model.File{
								{
									Path:   "example/path.go",
									Exists: true,
								},
							},
						},
					},
				},
			},
		}
		domainPlan := newPlan.Apply(model.Plan{})
		expectedDomainPlan := model.Plan{
			Revision:    1,
			Title:       lineage.NewRevisionedField("Plan Title"),
			Description: lineage.NewRevisionedField("Plan Description"),
			Steps: lineage.NewRevisionedField([]model.Step{
				{
					ID:      "S1",
					Title:   "Step Title",
					Summary: "Step Summary",
					ProposedChanges: []model.ProposedChange{
						{
							Description: "Proposed Change Description",
							Reason:      "Proposed Change Reason",
							Files: lineage.NewRevisionedField([]model.File{
								{
									Path:   "example/path.go",
									Exists: true,
								},
							}),
						},
					},
				},
			}),
		}
		testcmp.Compare(t, domainPlan, expectedDomainPlan)
	})

	t.Run("Agent Updates Existing Plan", func(t *testing.T) {
		t.Parallel()
		existingPlan := model.Plan{
			ID:          planID,
			Revision:    1,
			Title:       lineage.NewRevisionedField("Plan Title"),
			Description: lineage.NewRevisionedField("Plan Description"),
			Steps: lineage.NewRevisionedField([]model.Step{
				{
					ID:      "S1",
					Title:   "Step Title",
					Summary: "Step Summary",
					ProposedChanges: []model.ProposedChange{
						{
							Description: "Proposed Change Description",
							Reason:      "Proposed Change Reason",
							Files: lineage.NewRevisionedField([]model.File{
								{
									Path:   "example/path.go",
									Exists: true,
								},
							}),
						},
					},
				},
			}),
			CreatedAt: now,
		}

		// Convert domain plan to agent plan, for input into agent prompt.
		agentPlan := NewPlan(existingPlan)
		expectAgentPlanInput := Plan{
			Title:       "Plan Title",
			Description: "Plan Description",
			Steps: []Step{
				{
					ID:      "S1",
					Title:   "Step Title",
					Summary: "Step Summary",
					ProposedChanges: []ProposedChange{
						{
							Description: "Proposed Change Description",
							Reason:      "Proposed Change Reason",
							Files: []model.File{
								{
									Path:   "example/path.go",
									Exists: true,
								},
							},
						},
					},
				},
			},
		}
		testcmp.Compare(t, agentPlan, expectAgentPlanInput)

		agentPlan.DefinitionOfDone = []string{
			"foo",
			"bar",
		}
		updatedPlan := agentPlan.Apply(existingPlan)
		expectedUpdatedPlan := model.Plan{
			ID:          planID,
			Revision:    2,
			Title:       lineage.NewRevisionedField("Plan Title"),
			Description: lineage.NewRevisionedField("Plan Description"),
			DefinitionOfDone: lineage.RevisionedField[[]string]{
				Values: []lineage.RevisionedValue[[]string]{
					{Revision: 2, Value: []string{"foo", "bar"}},
				},
			},
			Steps: lineage.NewRevisionedField([]model.Step{
				{
					ID:      "S1",
					Title:   "Step Title",
					Summary: "Step Summary",
					ProposedChanges: []model.ProposedChange{
						{
							Description: "Proposed Change Description",
							Reason:      "Proposed Change Reason",
							Files: lineage.NewRevisionedField([]model.File{
								{
									Path:   "example/path.go",
									Exists: true,
								},
							}),
						},
					},
				},
			}),
			CreatedAt: now,
		}
		testcmp.Compare(t, updatedPlan, expectedUpdatedPlan)
	})
}
