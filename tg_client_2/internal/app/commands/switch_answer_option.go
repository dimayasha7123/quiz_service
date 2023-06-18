package commands

import (
	"context"
	"fmt"
	"github.com/dimayasha7123/quiz_service/tg_client_2/internal/app"
	"github.com/dimayasha7123/quiz_service/tg_client_2/internal/domain"
)

type SwitchReq struct {
	UserAnswersIDs app.UserAnswersIDs
}

type SwitchHandler struct {
	sessions domain.Sessions
}

func NewSwitchHandler(sessions domain.Sessions) SwitchHandler {
	return SwitchHandler{
		sessions: sessions,
	}
}

func (h SwitchHandler) Handle(ctx context.Context, req SwitchReq) error {
	err := h.sessions.SwitchAnswerForUser(ctx, req.UserAnswersIDs.UserID, req.UserAnswersIDs.AnswerID)
	if err != nil {
		return fmt.Errorf("can't switch answer with id = %v for user with id = %v: %v",
			req.UserAnswersIDs.UserID,
			req.UserAnswersIDs.AnswerID,
			err,
		)
	}
	return nil
}
