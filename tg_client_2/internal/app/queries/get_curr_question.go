package queries

import (
	"context"
	"fmt"
	"github.com/dimayasha7123/quiz_service/tg_client_2/internal/app"
	"github.com/dimayasha7123/quiz_service/tg_client_2/internal/domain"
)

type GetCurrQuestReq struct {
	UserID int64
}

type GetCurrQuestResp struct {
	QuestionInfo app.QuestionInfo
}

type getCurrQuestHandler struct {
	sessions domain.Sessions
}

func NewGetCurrQuestHandler(sessions domain.Sessions) getCurrQuestHandler {
	return getCurrQuestHandler{
		sessions: sessions,
	}
}

func (h getCurrQuestHandler) Handle(ctx context.Context, req GetCurrQuestReq) (GetCurrQuestResp, error) {
	questExists, question, err := h.sessions.GetCurrentQuestionForUser(ctx, req.UserID)
	if err != nil {
		return GetCurrQuestResp{}, fmt.Errorf("can't get current question for user with id = %v: %v", req.UserID, err)
	}
	if !questExists {
		return GetCurrQuestResp{
			QuestionInfo: app.QuestionInfo{
				Exist:    false,
				Question: app.Question{},
			},
		}, nil
	}

	answers := make(app.Answers, 0, len(question.Answers))
	for _, answer := range question.Answers {
		answers = append(answers, app.Answer{
			Title:  answer.Title,
			Picked: answer.Picked,
		})
	}
	questionApp := app.Question{
		Title:   question.Title,
		Answers: answers,
	}

	return GetCurrQuestResp{
		QuestionInfo: app.QuestionInfo{
			Exist:    true,
			Question: questionApp,
		},
	}, nil
}
