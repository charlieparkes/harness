package agentmodel

import (
	"errors"

	"github.com/charlieparkes/harness/internal/domain/model"
)

var ErrAnswerNotPopulated = errors.New("answer is not populated")

type Answer struct {
	QuestionID string `json:"question_id"`
	Answer     string `json:"answer"`
}

func NewAnswer(q model.Question) (Answer, error) {
	if q.Answer == nil {
		return Answer{}, ErrAnswerNotPopulated
	}
	return Answer{
		QuestionID: q.ID,
		Answer:     *q.Answer,
	}, nil
}
