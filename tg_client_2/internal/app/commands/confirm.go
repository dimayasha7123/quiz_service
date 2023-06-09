package commands

import (
	"context"
	"fmt"
	"github.com/dimayasha7123/quiz_service/server/pkg/api"
	"github.com/dimayasha7123/quiz_service/tg_client_2/internal/domain"
)

type ConfirmReq struct {
	UserID int64
}

type confirmHandler struct {
	sessions   domain.Sessions
	quizClient api.QuizServiceClient
}

func NewConfirmHandler(sessions domain.Sessions, quizClient api.QuizServiceClient) confirmHandler {
	return confirmHandler{
		sessions:   sessions,
		quizClient: quizClient,
	}
}

func (h confirmHandler) Handle(ctx context.Context, req ConfirmReq) error {
	err := h.sessions.ConfirmQuestionForUser(ctx, req.UserID)
	if err != nil {
		return fmt.Errorf("can't confirm question for user with id = %v: %v", req.UserID, err)
	}
	return nil
}
