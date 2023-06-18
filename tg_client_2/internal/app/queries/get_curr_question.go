package queries

import (
	"context"
	"fmt"
	"github.com/dimayasha7123/quiz_service/tg_client_2/internal/app"
	"github.com/dimayasha7123/quiz_service/tg_client_2/internal/domain"
)

type CurrentQuestReq struct {
	UserID int64
}

type CurrentQuestResp struct {
	QuestionInfo app.QuestionInfo
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
	questExists, question, err := h.sessions.CurrentQuestionForUser(ctx, req.UserID)
	if err != nil {
		return CurrentQuestResp{}, fmt.Errorf("can't get current question for user with id = %v: %v", req.UserID, err)
	}
	if !questExists {
		return CurrentQuestResp{
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

	return CurrentQuestResp{
		QuestionInfo: app.QuestionInfo{
			Exist:    true,
			Question: questionApp,
		},
	}, nil
}
