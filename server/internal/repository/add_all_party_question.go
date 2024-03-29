package repository

import (
	"context"
	"github.com/dimayasha7123/quiz_service/server/internal/models"
	"github.com/jackc/pgx/v4"
)

func (r *repository) AddAllPartyQuestion(ctx context.Context, quests []models.Question, partyID int64) error {

	rows := make([][]interface{}, len(quests))
	for i, q := range quests {
		rows[i] = []interface{}{q.ID, partyID, i}
	}

	_, err := r.pool.CopyFrom(
		ctx, pgx.Identifier{"party_question"},
		[]string{"question_id", "party_id", "quest_order_number"},
		pgx.CopyFromRows(rows),
	)

	return err
}
