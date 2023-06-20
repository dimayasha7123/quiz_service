package telegram_bot

import (
	"context"
	"fmt"
	"github.com/dimayasha7123/quiz_service/tg_client_2/internal/app/models"
	"github.com/dimayasha7123/quiz_service/tg_client_2/internal/app/queries"
	"github.com/dimayasha7123/quiz_service/tg_client_2/internal/ports/telegram_bot/callback"
	"github.com/dimayasha7123/quiz_service/tg_client_2/internal/ports/telegram_bot/commands"
	"strconv"
	"strings"
)

func (s *service) topHandle(ctx context.Context, req callbackQuery, callbackData callback.Data) error {
	if len(callbackData.Args) != 2 {
		return fmt.Errorf("must be exact 2 callbackData args, get %v", len(callbackData.Args))
	}

	quizID := callbackData.Args[0]
	quizTitle := callbackData.Args[1]

	quizIDInt, err := strconv.ParseInt(quizID, 10, 64)
	if err != nil {
		return fmt.Errorf("can't parse quizID = %v: %v", quizID, err)
	}

	topReq := queries.TopByQuizReq{
		UserQuizIDs: models.UserQuizIDs{
			UserID: req.From.ID,
			QuizID: quizIDInt,
		},
	}
	topResp, err := s.app.Queries.TopByQuizHandler.Handle(ctx, topReq)
	if err != nil {
		return fmt.Errorf("can't handle start req = %v: %v", topReq, err)
	}

	text := getTopText(topResp.Results, quizTitle)

	keyboard := [][]replyButton{
		{replyButton{
			Text:         "Back",
			CallbackData: callback.NewData(commands.Quiz, quizID, quizTitle).String(),
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

func getTopText(results models.Results, quizTitle string) string {
	sb := strings.Builder{}

	if len(results.TopResults) == 0 {
		sb.WriteString("No one has taken part in this quiz yet ((\nBut you can be the first!")
	} else {
		if quizTitle != "" {
			sb.WriteString(fmt.Sprintf("%s quiz results.\n\n", quizTitle))
		} else {
			sb.WriteString("Results.\n\n")
		}

		if results.UserResults.Place == 0 && results.UserResults.Points == 0 {
			sb.WriteString(fmt.Sprint("You didn't take part in this quiz ((\n\n"))
		} else {
			sb.WriteString("Your results:\n")
			sb.WriteString(fmt.Sprintf("Points: %d\n", results.UserResults.Points))
			sb.WriteString(fmt.Sprintf("Place: %d\n\n", results.UserResults.Place))
		}

		sb.WriteString("Top:\n")
		for place, rowResult := range results.TopResults {
			sb.WriteString(fmt.Sprintf("%d. @%s with %d p.\n", place+1, rowResult.Username, rowResult.Points))
		}
	}

	return sb.String()
}
