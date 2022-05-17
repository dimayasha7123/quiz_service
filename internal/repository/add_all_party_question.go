package repository

import (
	"context"
	"github.com/jackc/pgx/v4"
	"gitlab.ozon.dev/dimayasha7123/homework-2-dimayasha-7123/internal/models"
)

func (r *repository) AddAllPartyQuestion(ctx context.Context, quests []models.Question, partyID int64) error {

	rows := make([][]interface{}, len(quests))
	for i, q := range quests {
		rows[i] = make([]interface{}, 2)
		rows[i][0] = q.ID
		rows[i][1] = partyID
	}

	_, err := r.pool.CopyFrom(
		ctx, pgx.Identifier{"partyQuestion"},
		[]string{"question_id", "party_id"},
		pgx.CopyFromRows(rows),
	)

	return err
}
