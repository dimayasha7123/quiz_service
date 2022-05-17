package repository

import (
	"context"
	"gitlab.ozon.dev/dimayasha7123/homework-2-dimayasha-7123/internal/models"
)

func (r *repository) AddParty(ctx context.Context, party models.Party) (int64, error) {
	query := `
		insert into party (user_account_id, quiz_id)
		values ($1, $2) returning id;
	`

	var ID int64
	err := r.pool.QueryRow(ctx, query, party.UserAccountID, party.QuizID).Scan(&ID)

	return ID, err
}