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

func TestReviewStructFields(t *testing.T) {
	t.Parallel()
	testsize.Small(t)

	test.AssertStructFieldsEqual(t, Review{}, model.Review{},
		"CreatedAt",
		"UpdatedAt",
		"DeletedAt",
		"ID",
		"Revision",
		"TaskID",
	)
}

func TestReviewApplyRoundTrip(t *testing.T) {
	t.Parallel()
	testsize.Small(t)

	now := time.Date(2026, 9, 7, 12, 0, 0, 0, time.UTC)
	doc := testReviewDoc()
	review := doc.Apply(model.Review{ID: "review-1", TaskID: "task-1"}, now)

	require.Equal(t, int64(1), review.Revision)
	require.Equal(t, &now, review.UpdatedAt)
	require.Equal(t, "review-1", review.ID)
	require.Equal(t, "task-1", review.TaskID)

	got := NewReview(review)
	testcmp.Compare(t, got, doc)
}

func TestReviewApplyIdempotent(t *testing.T) {
	t.Parallel()
	testsize.Small(t)

	now := time.Date(2026, 9, 7, 12, 0, 0, 0, time.UTC)
	later := now.Add(time.Hour)
	doc := testReviewDoc()
	review := doc.Apply(model.Review{ID: "review-1", TaskID: "task-1"}, now)
	review = doc.Apply(review, later)

	require.Equal(t, int64(1), review.Revision)
	require.Equal(t, &now, review.UpdatedAt)

	require.Equal(t, []int64{1}, review.Decision.Revisions())
	require.Equal(t, []int64{1}, review.Observations.Revisions())
	require.Equal(t, []int64{1}, review.Verifications.Revisions())
	require.Equal(t, []int64{1}, review.Findings.Revisions())

	findings := review.Findings.Value()
	require.Equal(t, []int64{1}, findings[0].Evidence.Revisions())
	require.Equal(t, []int64{1}, findings[0].ProposedChanges.Revisions())

	evidence := findings[0].Evidence.Value()
	require.Equal(t, []int64{1}, evidence[0].Files.Revisions())
}

func TestReviewApplyPartialChange(t *testing.T) {
	t.Parallel()
	testsize.Small(t)

	now := time.Date(2026, 9, 7, 12, 0, 0, 0, time.UTC)
	later := now.Add(time.Hour)
	review := testReviewDoc().Apply(model.Review{ID: "review-1", TaskID: "task-1"}, now)

	line := int64(10)
	doc := testReviewDoc()
	doc.Decision = model.ReviewDecisionApproved
	doc.Findings[0].Evidence[0].Files = []model.File{
		{TargetID: "t1", Path: "a.go", Line: &line, Exists: true},
	}

	review = doc.Apply(review, later)
	require.Equal(t, int64(2), review.Revision)
	require.Equal(t, &later, review.UpdatedAt)

	require.Equal(t, []int64{1, 2}, review.Decision.Revisions())
	require.Equal(t, []int64{1}, review.Observations.Revisions())
	require.Equal(t, []int64{1, 2}, review.Findings.Revisions())

	findings := review.Findings.Value()
	require.Equal(t, []int64{1}, findings[0].Evidence.Revisions())
	require.Equal(t, []int64{1}, findings[1].Evidence.Revisions())

	evidence := findings[0].Evidence.Value()
	require.Equal(t, []int64{1}, evidence[0].Files.Revisions())
}

func testReviewDoc() Review {
	return Review{
		Decision:      model.ReviewDecisionChangesRequested,
		Observations:  []string{"observation 1", "observation 2"},
		Verifications: []string{"go test ./..."},
		Findings: []ReviewFinding{
			{
				ID:       "R1",
				Severity: model.ReviewFindingSeverityHigh,
				Category: model.ReviewFindingCategoryBug,
				Evidence: []ReviewFindingEvidence{
					{
						Description: "evidence 1 desc",
						Files:       []model.File{{TargetID: "t1", Path: "a.go", Exists: true}},
					},
				},
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
			},
			{
				ID:       "R2",
				Severity: model.ReviewFindingSeverityLow,
				Category: model.ReviewFindingCategoryBug,
				Evidence: []ReviewFindingEvidence{
					{
						Description: "evidence 2 desc",
						Files:       []model.File{{TargetID: "t2", Path: "b.go", Exists: true}},
					},
				},
				ProposedChanges: []ProposedChange{
					{
						Description: "change 2 desc",
						Reason:      "change 2 reason",
						Files:       []model.File{{TargetID: "t2", Path: "b.go", Exists: true}},
					},
				},
				ProposedTests: []ProposedTest{
					{
						Description: "test change 2 desc",
						Reason:      "test change 2 reason",
						Files:       []model.File{{TargetID: "t2", Path: "b_test.go", Exists: false}},
						TestCases:   []string{"case 1"},
					},
				},
			},
		},
	}
}
