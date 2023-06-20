package telegram_bot

import (
	"context"
	"fmt"
	"github.com/dimayasha7123/quiz_service/tg_client_2/internal/ports/telegram_bot/callback"
	"github.com/dimayasha7123/quiz_service/tg_client_2/internal/ports/telegram_bot/commands"
)

const quizzesText = "Choose quiz:"

func (s *service) quizzesHandle(ctx context.Context, req message) error {
	keyboard, err := s.getQuizzesKeyboard(ctx)
	if err != nil {
		return fmt.Errorf("can't get quizzes keyboard for: %v", err)
	}

	err = s.sendMessage(ctx, req.From.ID, quizzesText, &replyMarkup{InlineKeyboard: keyboard})
	if err != nil {
		return fmt.Errorf("can't send message to chat with id = %v: %v", req.From.ID, err)
	}

	return nil
}

func (s *service) backToQuizListHandle(ctx context.Context, req callbackQuery, callbackData callback.Data) error {
	keyboard, err := s.getQuizzesKeyboard(ctx)
	if err != nil {
		return fmt.Errorf("can't get quizzes keyboard for: %v", err)
	}

	err = s.editMessageAndReplyMarkup(ctx, req.Message.Chat.ID, req.Message.MessageID, quizzesText, replyMarkup{InlineKeyboard: keyboard})
	if err != nil {
		return fmt.Errorf(
			"can't edit message and reply markup with chatID = %d, messID = %d, text = %s and keyboard = %v: %v",
			req.Message.Chat.ID,
			req.Message.MessageID,
			quizzesText,
			keyboard,
			err,
		)
	}

	return nil
}

func (s *service) getQuizzesKeyboard(ctx context.Context) ([][]replyButton, error) {
	quizzes, err := s.app.Queries.QuizzesHandler.Handle(ctx)
	if err != nil {
		return nil, fmt.Errorf("can't handle quizzes req: %v", err)
	}
	if len(quizzes.QuizList) == 0 {
		return nil, fmt.Errorf("get 0 quizzes")
	}

	keyboard := make([][]replyButton, 0, len(quizzes.QuizList))
	for _, quiz := range quizzes.QuizList {
		keyboard = append(keyboard, []replyButton{{
			Text:         quiz.Title,
			CallbackData: callback.NewData(commands.Quiz, fmt.Sprint(quiz.ID), quiz.Title).String(),
		}})
	}

	return keyboard, nil
}
