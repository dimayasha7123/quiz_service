package app

import (
	"context"
	"errors"
	"strconv"

	"github.com/dimayasha7123/quiz_service/tg_client/internal/models"
)

func (b *bclient) pickUnpickHandler(ctx context.Context, update models.Update, commandArgs []string, value bool) (string, error) {

	if len(commandArgs) != 1 {
		return "", errors.New("bad arguments")
	}

	b.users.RLock()
	user, ok := b.users.M[update.Message.From.ID]
	b.users.RUnlock()

	if !ok {
		return "", errors.New("user not found")
	}
	if user.State != 1 {
		return "", errors.New("invalid state")
	}

	ansNum, err := strconv.ParseInt(commandArgs[0], 10, 64)
	if err != nil {
		return "", errors.New("bad arguments")
	}

	user.Questions[user.CurrentQuestion].AnswerOptions[ansNum-1].Picked = value
	b.users.SafeUpdateUser(user)

	text, ok := user.GetQuestion(user.CurrentQuestion)
	if !ok {
		return "", errors.New("quiz has no quiestions")
	}

	return text, nil
}
