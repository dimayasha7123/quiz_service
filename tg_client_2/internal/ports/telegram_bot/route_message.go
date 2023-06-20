package telegram_bot

import (
	"context"
	"fmt"
	"github.com/dimayasha7123/quiz_service/utils/logger"
)

func (s *service) routeMessage(ctx context.Context, req message) error {
	command := req.Text

	// TODO: make normal logging!!!
	logger.Log.Infof("get text request %q", command)

	var err error
	switch command {
	case "/start":
		err = s.startHandle(ctx, req)
	case "/help":
		err = s.helpHandle(ctx, req)
	case "/quizzes":
		err = s.quizzesHandle(ctx, req)
	case "/break":
		err = s.breakHandle(ctx, req)
	case "/canon":
		err = s.canonHandle(ctx, req)
	default:
		err = fmt.Errorf("no such command")
	}

	if err != nil {
		return fmt.Errorf("can't handle %q command: %v", command, err)
	}

	return nil
}
