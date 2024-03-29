package queries

import (
	"context"
	"fmt"
	"github.com/dimayasha7123/quiz_service/tg_client_2/internal/app/models"
	"github.com/dimayasha7123/quiz_service/tg_client_2/internal/domain"
)

type CurrentQuestReq struct {
	UserID int64
}

type CurrentQuestResp struct {
	QuestionInfo models.QuestionInfo
}

type CurrentQuestHandler struct {
	sessions domain.Sessions
}

func NewCurrentQuestHandler(sessions domain.Sessions) CurrentQuestHandler {
	return CurrentQuestHandler{
		sessions: sessions,
	}
}

func (h CurrentQuestHandler) Handle(ctx context.Context, req CurrentQuestReq) (CurrentQuestResp, error) {
	questExists, number, question, err := h.sessions.CurrentQuestionForUser(ctx, req.UserID)
	if err != nil {
		return CurrentQuestResp{}, fmt.Errorf("can't get current question for user with id = %v: %v", req.UserID, err)
	}
	if !questExists {
		return CurrentQuestResp{
			QuestionInfo: models.QuestionInfo{
				Exist:    false,
				Number:   number,
				Question: models.Question{},
			},
		}, nil
	}

	answers := make(models.Answers, 0, len(question.Answers))
	for _, answer := range question.Answers {
		answers = append(answers, models.Answer{
			Title:  answer.Title,
			Picked: answer.Picked,
		})
	}
	questionApp := models.Question{
		Title:   question.Title,
		Answers: answers,
	}

	return CurrentQuestResp{
		QuestionInfo: models.QuestionInfo{
			Exist:    true,
			Number:   number,
			Question: questionApp,
		},
	}, nil
}
