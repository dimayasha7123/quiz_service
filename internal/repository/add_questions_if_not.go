package repository

import (
	"context"
	"github.com/jackc/pgx/v4"
	"gitlab.ozon.dev/dimayasha7123/homework-2-dimayasha-7123/internal/models"
)

func (r *repository) AddQuestionsIfNot(ctx context.Context, quests []models.Question, quizID int64) (int32, error) {

	query := `
		select id
		from question
		where title != ($1);
	`

	titles := make([]string, len(quests))

	for i, q := range quests {
		titles[i] = q.Title
	}

	rows, err := r.pool.Query(ctx, query, titles)
	if err != nil {
		return -1, err
	}

	notInIDs := make([]int64, 0, len(quests))

	for rows.Next() {
		var id int64
		err = rows.Scan(&id)
		if err != nil {
			return -1, err
		}
		notInIDs = append(notInIDs, id)
	}

	questsToAdd := make([][]interface{}, 0, len(quests))
	answersToAdd := make([][]interface{}, 0, len(quests)*4)

	for _, qst := range quests {

		notIn := false
		for _, nID := range notInIDs {
			if qst.ID == nID {
				notIn = true
				break
			}
		}
		if notIn {
			continue
		}

		questsToAdd = append(questsToAdd, []interface{}{qst.Title, quizID})
		for _, ans := range qst.Answers {
			answersToAdd = append(answersToAdd, []interface{}{qst.ID, ans.Title, ans.Correct})
		}
	}

	added, err := r.pool.CopyFrom(
		ctx, pgx.Identifier{"question"},
		[]string{"title", "quiz_id"},
		pgx.CopyFromRows(questsToAdd),
	)
	if err != nil {
		return -1, err
	}

	_, err = r.pool.CopyFrom(
		ctx, pgx.Identifier{"answer"},
		[]string{"question_id", "title", "correct"},
		pgx.CopyFromRows(answersToAdd),
	)
	if err != nil {
		return -1, err
	}

	return int32(added), nil
}
