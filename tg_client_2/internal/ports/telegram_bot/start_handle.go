package telegram_bot

import (
	"context"
	"fmt"
	"github.com/dimayasha7123/quiz_service/tg_client_2/internal/app/models"
	"github.com/dimayasha7123/quiz_service/tg_client_2/internal/app/queries"
	"strings"
)

func (s *service) startHandle(ctx context.Context, req message) error {
	startReq := queries.StartReq{
		UserInfo: models.UserInfo{
			UserID:   req.From.ID,
			UserName: req.From.Username,
		}}
	startResp, err := s.app.Queries.StartHandler.Handle(ctx, startReq)
	if err != nil {
		return fmt.Errorf("can't handle start req = %v: %v", startReq, err)
	}

	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("Welcome, %s.", req.From.Username))
	if !startResp.NewUser {
		sb.WriteString(" Again.")
	}
	sb.WriteString("\n\n")
	sb.WriteString("Quick start: /quizzes\n")
	sb.WriteString("Learn more: /help\n")

	err = s.sendMessage(ctx, req.From.ID, sb.String(), nil)
	if err != nil {
		return fmt.Errorf("can't send message to chat with id = %v: %v", req.From.ID, err)
	}

	return nil
}
