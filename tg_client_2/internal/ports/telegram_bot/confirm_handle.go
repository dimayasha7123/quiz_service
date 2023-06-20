package telegram_bot

import (
	"context"
	"fmt"
	"github.com/dimayasha7123/quiz_service/tg_client_2/internal/app/commands"
	"github.com/dimayasha7123/quiz_service/tg_client_2/internal/app/queries"
	"github.com/dimayasha7123/quiz_service/tg_client_2/internal/ports/telegram_bot/callback"
	commands2 "github.com/dimayasha7123/quiz_service/tg_client_2/internal/ports/telegram_bot/commands"
)

func (s *service) confirmHandle(ctx context.Context, req callbackQuery, callbackData callback.Data) error {
	confirmReq := commands.ConfirmReq{
		UserID: req.From.ID,
	}
	err := s.app.Commands.ConfirmHandler.Handle(ctx, confirmReq)
	if err != nil {
		return fmt.Errorf("can't handle confirm req = %v: %v", confirmReq, err)
	}

	currQuestReq := queries.CurrentQuestReq{
		UserID: req.From.ID,
	}
	currQuestResp, err := s.app.Queries.CurrQuestHandler.Handle(ctx, currQuestReq)
	if err != nil {
		return fmt.Errorf("can't handle current quest req = %v: %v", currQuestReq, err)
	}

	if !currQuestResp.QuestionInfo.Exist {
		endReq := queries.EnqQuizReq{
			UserID: req.From.ID,
		}
		endResp, err := s.app.Queries.EnqQuizHandler.Handle(ctx, endReq)
		if err != nil {
			return fmt.Errorf("can't handle end req = %v: %v", endReq, err)
		}

		// TODO: can't get quizTitle here, add ifchik
		text := getTopText(endResp.Results, "")

		keyboard := [][]replyButton{
			{replyButton{
				Text:         "Back to quizzes",
				CallbackData: callback.NewData(commands2.BackToQuizList).String(),
			}},
		}

		err = s.editMessageAndReplyMarkup(ctx, req.Message.Chat.ID, req.Message.MessageID, text, replyMarkup{InlineKeyboard: keyboard})
		if err != nil {
			return fmt.Errorf(
				"can't edit message and reply markup with chatID = %d, messID = %d, text = %s and keyboard = %v: %v",
				req.Message.Chat.ID,
				req.Message.MessageID,
				text,
				keyboard,
				err,
			)
		}

		return nil
	}

	err = s.editCurrentQuestion(ctx, req, currQuestResp)
	if err != nil {
		return fmt.Errorf("can't edit current question for req = %v: %v", req, err)
	}

	return nil
}
