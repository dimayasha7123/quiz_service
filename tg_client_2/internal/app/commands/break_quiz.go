package commands

import (
	"context"
	"fmt"
	"github.com/dimayasha7123/quiz_service/tg_client_2/internal/domain"
)

type BreakReq struct {
	UserID int64
}

type BreakHandler struct {
	sessions domain.Sessions
}

func NewBreakHandler(sessions domain.Sessions) BreakHandler {
	return BreakHandler{
		sessions: sessions,
	}
}

func (h BreakHandler) Handle(ctx context.Context, req BreakReq) error {
	err := h.sessions.BreakSessionForUser(ctx, req.UserID)
	if err != nil {
		return fmt.Errorf(
			"can't break session for user with id = %v: %v",
			req.UserID,
			err,
		)
	}
	return nil
}
