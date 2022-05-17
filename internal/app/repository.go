package app

import (
	"context"
	"gitlab.ozon.dev/dimayasha7123/homework-2-dimayasha-7123/internal/models"
)

type Repository interface {
	AddUser(context.Context, string) (int64, error)
	GetQuizList(context.Context) ([]models.Quiz, error)
	AddParty(context.Context, models.Party) (int64, error)
}
