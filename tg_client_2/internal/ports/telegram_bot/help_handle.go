package telegram_bot

import (
	"context"
	"fmt"
	"strings"
)

func (s *service) helpHandle(ctx context.Context, req message) error {
	sb := strings.Builder{}
	sb.WriteString("This bot allows you to pass quizzes on IT topics.\n")
	sb.WriteString("Choose a topic, answer questions, you can choose several answers, and at the end you get points and a rating.\n\n")
	sb.WriteString("To begin use command /quizzes.\n")
	sb.WriteString("If you have started a quiz and don't want to finish it, then use command /break.\n")

	err := s.sendMessage(ctx, req.From.ID, sb.String(), nil)
	if err != nil {
		return fmt.Errorf("can't send message to chat with id = %v: %v", req.From.ID, err)
	}

	return nil
}
