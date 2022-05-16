package repository

import (
	"context"
	"gitlab.ozon.dev/dimayasha7123/homework-2-dimayasha-7123/internal/models"
)

func (r *repository) AddParticipation(ctx context.Context, party models.Participation) (int64, error) {
	query := `
		insert into participation (user_account_id, quiz_id)
		values ($1, $2) returning id;
	`

	var ID int64
	err := r.pool.QueryRow(ctx, query, party.UserAccountID, party.QuizID).Scan(&ID)

	return ID, err
}
