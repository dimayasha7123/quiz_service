package telegram_bot

import (
	"context"
	"fmt"
	"github.com/dimayasha7123/quiz_service/tg_client_2/internal/app/commands"
	"github.com/dimayasha7123/quiz_service/tg_client_2/internal/app/models"
	"github.com/dimayasha7123/quiz_service/tg_client_2/internal/app/queries"
	"github.com/dimayasha7123/quiz_service/tg_client_2/internal/ports/telegram_bot/callback"
	"strconv"
)

func (s *service) switchAnswerHandle(ctx context.Context, req callbackQuery, callbackData callback.Data) error {
	if len(callbackData.Args) != 1 {
		return fmt.Errorf("must be exact 1 callbackData args, get %v", len(callbackData.Args))
	}
	answerID := callbackData.Args[0]
	answerIDInt, err := strconv.ParseInt(answerID, 10, 64)
	if err != nil {
		return fmt.Errorf("can't parse answerID = %v: %v", answerID, err)
	}

	switchReq := commands.SwitchReq{
		UserAnswersIDs: models.UserAnswersIDs{
			UserID:   req.From.ID,
			AnswerID: answerIDInt,
		},
	}
	err = s.app.Commands.SwitchHandler.Handle(ctx, switchReq)
	if err != nil {
		return fmt.Errorf("can't handle switch asnwer req = %v: %v", switchReq, err)
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
