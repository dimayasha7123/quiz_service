package telegram_bot

import (
	"context"
	"fmt"
	"github.com/dimayasha7123/quiz_service/tg_client_2/internal/app/commands"
)

func (s *service) breakHandle(ctx context.Context, req message) error {
	breakReq := commands.BreakReq{
		UserID: req.From.ID,
	}
	err := s.app.Commands.BreakHandler.Handle(ctx, breakReq)
	if err != nil {
		return fmt.Errorf("can't handle break req = %v: %v", breakReq, err)
	}

	err = s.sendMessage(ctx, req.From.ID, "Out of quiz! ((", nil)
	if err != nil {
		return fmt.Errorf("can't send message to chat with id = %v: %v", req.From.ID, err)
	}

	return nil
}
