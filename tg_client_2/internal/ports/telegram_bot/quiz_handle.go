package telegram_bot

import (
	"context"
	"fmt"
	"github.com/dimayasha7123/quiz_service/tg_client_2/internal/ports/telegram_bot/callback"
	"github.com/dimayasha7123/quiz_service/tg_client_2/internal/ports/telegram_bot/commands"
)

func (s *service) quizHandle(ctx context.Context, req callbackQuery, callbackData callback.Data) error {
	if len(callbackData.Args) != 2 {
		return fmt.Errorf("must be exact 2 callbackData args, get %v", len(callbackData.Args))
	}

	quizID := callbackData.Args[0]
	quizTitle := callbackData.Args[1]

	keyboard := [][]replyButton{
		{replyButton{
			Text:         "Take part",
			CallbackData: callback.NewData(commands.TakeParty, quizID).String(),
		}},
		{replyButton{
			Text:         "Top",
			CallbackData: callback.NewData(commands.LookTop, quizID, quizTitle).String(),
		}},
		{replyButton{
			Text:         "Back",
			CallbackData: callback.NewData(commands.BackToQuizList).String(),
		}},
	}

	err := s.editMessageAndReplyMarkup(ctx, req.Message.Chat.ID, req.Message.MessageID, quizTitle, replyMarkup{InlineKeyboard: keyboard})
	if err != nil {
		return fmt.Errorf(
			"can't edit message and reply markup with chatID = %d, messID = %d, text = %s and keyboard = %v: %v",
			req.Message.Chat.ID,
			req.Message.MessageID,
			quizTitle,
			keyboard,
			err,
		)
	}

	return nil
}

// TODO: для старта квиза
//startQuizReq := queries.StartQuizReq{
//	UserQuizIDs: models.UserQuizIDs{
//		UserID: req.From.ID,
//		QuizID: quizID,
//	},
//}
//startQuizResp, err := s.app.Queries.StartQuizHandler.Handle(ctx, startQuizReq)
//if err != nil {
//	return fmt.Errorf("can't handle start quiz req = %v: %v", startQuizReq, err)
//}
