package agentmodel

import (
	"testing"
	"time"

	"github.com/charlieparkes/go-testcmp"
	"github.com/charlieparkes/go-testsize"
	"github.com/charlieparkes/harness/internal/domain/model"
	"github.com/charlieparkes/harness/test"
	"github.com/stretchr/testify/require"
)

func TestPlanStructFields(t *testing.T) {
	t.Parallel()
	testsize.Small(t)

	test.AssertStructFieldsEqual(t, Plan{}, model.Plan{},
		"CreatedAt",
		"UpdatedAt",
		"DeletedAt",
		"ID",
		"Revision",
		"TaskID",
		"Status",
	)
}

func TestPlanApplyRoundTrip(t *testing.T) {
	t.Parallel()
	testsize.Small(t)

	now := time.Date(2026, 9, 7, 12, 0, 0, 0, time.UTC)
	doc := testPlanDoc()
	plan := doc.Apply(model.Plan{ID: "plan-1", TaskID: "task-1"}, now)

	require.Equal(t, int64(1), plan.Revision)
	require.Equal(t, &now, plan.UpdatedAt)
	require.Equal(t, "plan-1", plan.ID)
	require.Equal(t, "task-1", plan.TaskID)

	got := NewPlan(plan)
	testcmp.Compare(t, got, doc)
}

func TestPlanApplyIdempotent(t *testing.T) {
	t.Parallel()
	testsize.Small(t)

	now := time.Date(2026, 9, 7, 12, 0, 0, 0, time.UTC)
	later := now.Add(time.Hour)
	doc := testPlanDoc()
	plan := doc.Apply(model.Plan{ID: "plan-1", TaskID: "task-1"}, now)
	plan = doc.Apply(plan, later)

	require.Equal(t, int64(1), plan.Revision)
	require.Equal(t, &now, plan.UpdatedAt)

	require.Equal(t, []int64{1}, plan.Title.Revisions())
	require.Equal(t, []int64{1}, plan.Description.Revisions())
	require.Equal(t, []int64{1}, plan.Risks.Revisions())
	require.Equal(t, []int64{1}, plan.DefinitionOfDone.Revisions())
	require.Equal(t, []int64{1}, plan.Steps.Revisions())

	steps, _ := plan.Steps.Value()
	require.Equal(t, []int64{1}, steps[0].ProposedChanges[0].Files.Revisions())
}

func TestPlanApplyPartialChange(t *testing.T) {
	t.Parallel()
	testsize.Small(t)

	now := time.Date(2026, 9, 7, 12, 0, 0, 0, time.UTC)
	later := now.Add(time.Hour)
	plan := testPlanDoc().Apply(model.Plan{ID: "plan-1", TaskID: "task-1"}, now)

	line := int64(10)
	doc := testPlanDoc()
	doc.Title = "Updated title"
	doc.Steps[0].ProposedChanges[0].Files = []model.File{
		{TargetID: "t1", Path: "a.go", Line: &line, Exists: true},
	}

	plan = doc.Apply(plan, later)
	require.Equal(t, int64(2), plan.Revision)
	require.Equal(t, &later, plan.UpdatedAt)

	require.Equal(t, []int64{1, 2}, plan.Title.Revisions())
	require.Equal(t, []int64{1}, plan.Description.Revisions())
	require.Equal(t, []int64{1, 2}, plan.Steps.Revisions())

	steps, _ := plan.Steps.Value()
	require.Equal(t, []int64{1, 2}, steps[0].ProposedChanges[0].Files.Revisions())
	require.Equal(t, []int64{1}, steps[1].ProposedChanges[0].Files.Revisions())
}

func testPlanDoc() Plan {
	return Plan{
		Title:            "Add feature",
		Description:      "Detailed description",
		DefinitionOfDone: []string{"done 1", "done 2"},
		Risks:            []string{"plan risk"},
		Steps: []Step{
			{
				ID:      "S1",
				Title:   "First step",
				Summary: "Summary 1",
				ProposedChanges: []ProposedChange{
					{
						Description: "change desc",
						Reason:      "change reason",
						Files:       []model.File{{TargetID: "t1", Path: "a.go", Exists: true}},
					},
				},
				ProposedTests: []ProposedTest{
					{
						Description: "test change desc",
						Reason:      "test change reason",
						Files:       []model.File{{TargetID: "t1", Path: "a_test.go", Exists: false}},
						TestCases:   []string{"case 1", "case 2"},
					},
				},
				Verifications: []string{"go test ./..."},
				Risks:         []string{"step risk"},
				Dependencies:  []string{},
			},
			{
				ID:      "S2",
				Title:   "Second step",
				Summary: "Summary 2",
				ProposedChanges: []ProposedChange{
					{
						Description: "change 2 desc",
						Reason:      "change 2 reason",
						Files:       []model.File{{TargetID: "t2", Path: "b.go", Exists: true}},
					},
				},
				ProposedTests: []ProposedTest{},
				Verifications: []string{"make"},
				Risks:         []string{},
				Dependencies:  []string{"S1"},
			},
		},
	}
}
