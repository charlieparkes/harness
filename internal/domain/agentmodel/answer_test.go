package agentmodel

import (
	"testing"

	"github.com/charlieparkes/go-testcmp"
	"github.com/charlieparkes/go-testsize"
	"github.com/charlieparkes/harness/internal/domain/model"
	"github.com/stretchr/testify/require"
)

func TestNewAnswer(t *testing.T) {
	t.Parallel()
	testsize.Small(t)

	answer := "Use SQLite"
	got, err := NewAnswer(model.Question{
		ID:       "Q1",
		Question: "Which store should Harness use?",
		Answer:   &answer,
	})
	require.NoError(t, err)
	testcmp.Compare(t, got, Answer{
		QuestionID: "Q1",
		Answer:     "Use SQLite",
	})
}

func TestNewAnswerEmptyString(t *testing.T) {
	t.Parallel()
	testsize.Small(t)

	answer := ""
	got, err := NewAnswer(model.Question{
		ID:     "Q1",
		Answer: &answer,
	})
	require.NoError(t, err)
	testcmp.Compare(t, got, Answer{
		QuestionID: "Q1",
		Answer:     "",
	})
}

func TestNewAnswerNotPopulated(t *testing.T) {
	t.Parallel()
	testsize.Small(t)

	_, err := NewAnswer(model.Question{ID: "Q1"})
	require.ErrorIs(t, err, ErrAnswerNotPopulated)
}
