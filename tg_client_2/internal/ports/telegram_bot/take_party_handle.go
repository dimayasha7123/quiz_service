package telegram_bot

import (
	"context"
	"fmt"
	"github.com/dimayasha7123/quiz_service/tg_client_2/internal/app/commands"
	"github.com/dimayasha7123/quiz_service/tg_client_2/internal/app/models"
	"github.com/dimayasha7123/quiz_service/tg_client_2/internal/app/queries"
	"github.com/dimayasha7123/quiz_service/tg_client_2/internal/ports/telegram_bot/callback"
	commands2 "github.com/dimayasha7123/quiz_service/tg_client_2/internal/ports/telegram_bot/commands"
	"strconv"
	"strings"
)

const (
	unpickedSymbol = "🔴"
	pickedSymbol   = "🟢"
)

func (s *service) takePartyHandle(ctx context.Context, req callbackQuery, callbackData callback.Data) error {
	if len(callbackData.Args) != 1 {
		return fmt.Errorf("must be exact 1 callbackData args, get %v", len(callbackData.Args))
	}
	quizID := callbackData.Args[0]
	quizIDInt, err := strconv.ParseInt(quizID, 10, 64)
	if err != nil {
		return fmt.Errorf("can't parse quizID = %v: %v", quizID, err)
	}

	startQuizReq := commands.StartQuizReq{
		UserQuizIDs: models.UserQuizIDs{
			UserID: req.From.ID,
			QuizID: quizIDInt,
		},
	}

	err = s.app.Commands.StartQuizHandler.Handle(ctx, startQuizReq)
	if err != nil {
		return fmt.Errorf("can't handle start quiz req = %v: %v", startQuizReq, err)
	}

	currQuestReq := queries.CurrentQuestReq{
		UserID: req.From.ID,
	}
	currQuestResp, err := s.app.Queries.CurrQuestHandler.Handle(ctx, currQuestReq)
	if err != nil {
		return fmt.Errorf("can't handle current quest req = %v: %v", currQuestReq, err)
	}
	if !currQuestResp.QuestionInfo.Exist {
		return fmt.Errorf("quest not exist (but must exist)")
	}

	err = s.editCurrentQuestion(ctx, req, currQuestResp)
	if err != nil {
		return fmt.Errorf("can't edit current question for req = %v: %v", req, err)
	}

	return nil
}

func (s *service) editCurrentQuestion(ctx context.Context, req callbackQuery, currQuestResp queries.CurrentQuestResp) error {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("%d. %s\n\n", currQuestResp.QuestionInfo.Number, currQuestResp.QuestionInfo.Question.Title))
	sb.WriteString("Answers:\n")
	answers := currQuestResp.QuestionInfo.Question.Answers
	for i, ans := range answers {
		sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, ans.Title))
	}
	questionText := sb.String()

	keyboard := make([][]replyButton, 0, len(answers))
	for i, answer := range answers {
		symbol := unpickedSymbol
		if answer.Picked {
			symbol = pickedSymbol
		}
		keyboard = append(keyboard, []replyButton{{
			Text:         fmt.Sprintf("%s %s", symbol, answer.Title),
			CallbackData: callback.NewData(commands2.SwitchAnswer, fmt.Sprint(i)).String(),
		}})
	}
	keyboard = append(keyboard, []replyButton{{
		Text:         "Confirm",
		CallbackData: callback.NewData(commands2.Confirm).String(),
	}})

	err := s.editMessageAndReplyMarkup(ctx, req.Message.Chat.ID, req.Message.MessageID, questionText, replyMarkup{InlineKeyboard: keyboard})
	if err != nil {
		return fmt.Errorf(
			"can't edit message and reply markup with chatID = %d, messID = %d, text = %s and keyboard = %v: %v",
			req.Message.Chat.ID,
			req.Message.MessageID,
			questionText,
			keyboard,
			err,
		)
	}

	return nil
}
