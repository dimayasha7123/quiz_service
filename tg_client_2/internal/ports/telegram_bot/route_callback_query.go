package telegram_bot

import (
	"context"
	"fmt"
	"github.com/dimayasha7123/quiz_service/tg_client_2/internal/ports/telegram_bot/callback"
	"github.com/dimayasha7123/quiz_service/tg_client_2/internal/ports/telegram_bot/commands"
	"github.com/dimayasha7123/quiz_service/utils/logger"
)

type NotImplementedError struct{}

func (e NotImplementedError) Error() string {
	return "not implemented yet"
}

func (s *service) routeCallbackQuery(ctx context.Context, req callbackQuery) error {

	// TODO: make normal logging!!!
	logger.Log.Infof("get callback request %q", req.Data)

	callbackData, callbackErr := callback.NewDataFromString(req.Data)
	if callbackErr != nil {
		return fmt.Errorf("can't create callbackData from string = %v: %v", req.Data, callbackErr)
	}

	var err error
	switch callbackData.NextCommand {
	case commands.Quiz:
		err = s.quizHandle(ctx, req, callbackData)
	case commands.TakeParty:
		err = s.takePartyHandle(ctx, req, callbackData)
	case commands.LookTop:
		err = s.topHandle(ctx, req, callbackData)
	case commands.BackToQuizList:
		err = s.backToQuizListHandle(ctx, req, callbackData)
	case commands.SwitchAnswer:
		err = s.switchAnswerHandle(ctx, req, callbackData)
	case commands.Confirm:
		err = s.confirmHandle(ctx, req, callbackData)
	default:
		err = fmt.Errorf("unknown callback command")
	}

	if err != nil {
		return fmt.Errorf("can't handle %q callback: %v", callbackData.NextCommand.String(), err)
	}

	return nil
}
