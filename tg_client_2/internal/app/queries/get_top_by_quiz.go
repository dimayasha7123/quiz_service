package queries

import (
	"context"
	"fmt"
	"github.com/dimayasha7123/quiz_service/server/pkg/api"
	"github.com/dimayasha7123/quiz_service/tg_client_2/internal/app/models"
	"github.com/dimayasha7123/quiz_service/tg_client_2/internal/domain"
)

type TopByQuizReq struct {
	UserQuizIDs models.UserQuizIDs
}

type TopByQuizResp struct {
	Results models.Results
}

type TopByQuizHandler struct {
	sessions   domain.Sessions
	quizClient api.QuizServiceClient
}

func NewTopByQuizHandler(sessions domain.Sessions, quizClient api.QuizServiceClient) TopByQuizHandler {
	return TopByQuizHandler{
		sessions:   sessions,
		quizClient: quizClient,
	}
}

func (h TopByQuizHandler) Handle(ctx context.Context, req TopByQuizReq) (TopByQuizResp, error) {
	user, err := h.sessions.UserByID(ctx, req.UserQuizIDs.UserID)
	if err != nil {
		return TopByQuizResp{}, fmt.Errorf("can't get user from sessions: %v", err)
	}

	qcResp, err := h.quizClient.GetQuizTop(ctx, &api.QuizUserInfo{
		UserID: user.QuizServiceID,
		QuizID: req.UserQuizIDs.QuizID,
	})

	results := convertResultsFromApiToApp(ctx, h.sessions, qcResp)

	return TopByQuizResp{Results: results}, nil
}
