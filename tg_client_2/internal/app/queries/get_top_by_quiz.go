package queries

import (
	"context"
	"fmt"
	"github.com/dimayasha7123/quiz_service/server/pkg/api"
	"github.com/dimayasha7123/quiz_service/tg_client_2/internal/app"
	"github.com/dimayasha7123/quiz_service/tg_client_2/internal/domain"
)

type GetTopByQuizReq struct {
	UserQuizIDs app.UserQuizIDs
}

type GetTopByQuizResp struct {
	Results app.Results
}

type getTopByQuizHandler struct {
	sessions   domain.Sessions
	quizClient api.QuizServiceClient
}

func NewGetTopByQuizHandler(sessions domain.Sessions, quizClient api.QuizServiceClient) getTopByQuizHandler {
	return getTopByQuizHandler{
		sessions:   sessions,
		quizClient: quizClient,
	}
}

func (h getTopByQuizHandler) Handle(ctx context.Context, req GetTopByQuizReq) (GetTopByQuizResp, error) {
	user, err := h.sessions.GetUserByID(ctx, req.UserQuizIDs.UserID)
	if err != nil {
		return GetTopByQuizResp{}, fmt.Errorf("can't get user from sessions: %v", err)
	}

	qcResp, err := h.quizClient.GetQuizTop(ctx, &api.QuizUserInfo{
		UserID: user.QuizServiceID,
		QuizID: req.UserQuizIDs.QuizID,
	})

	results := convertResultsFromApiToApp(ctx, h.sessions, qcResp)

	return GetTopByQuizResp{Results: results}, nil
}
