package repository

import (
	"context"
	"github.com/dimayasha7123/quiz_service/server/internal/models"
	"github.com/jackc/pgx/v4"
)

func (r *repository) AddQuestionsIfNot(ctx context.Context, questsPointer *[]models.Question, quizID int64) (int32, error) {

	query := `
		select id, title
		from question
		where title = any ($1);
	`

	qCount := len(*questsPointer)
	quests := *questsPointer

	titles := make([]string, qCount)

	for i, q := range quests {
		titles[i] = q.Title
	}
	rows, err := r.pool.Query(ctx, query, titles)
	if err != nil {
		return -1, err
	}

	inQuizes := make([]models.Quiz, 0, qCount)

	for rows.Next() {
		var id int64
		var t string
		err = rows.Scan(&id, &t)
		if err != nil {
			return -1, err
		}
		inQuizes = append(inQuizes, models.Quiz{ID: id, Title: t})
	}

	answersToAdd := make([][]interface{}, 0, qCount*4)

	added := 0

	for i := range quests {

		next := false
		for _, inQuiz := range inQuizes {
			if quests[i].Title == inQuiz.Title {
				quests[i].ID = inQuiz.ID
				next = true
				break
			}
		}
		if next {
			continue
		}

		questQuery := `
			insert into question (title, quiz_id) 
			values ($1, $2) returning id;
		`

		err = r.pool.QueryRow(ctx, questQuery, quests[i].Title, quizID).Scan(&quests[i].ID)
		if err != nil {
			return -1, err
		}
		added++

		for j, ans := range quests[i].Answers {
			answersToAdd = append(answersToAdd, []interface{}{quests[i].ID, ans.Title, ans.Correct, j})
		}
	}

	_, err = r.pool.CopyFrom(
		ctx, pgx.Identifier{"answer"},
		[]string{"question_id", "title", "correct", "order_number"},
		pgx.CopyFromRows(answersToAdd),
	)
	if err != nil {
		return -1, err
	}

	return int32(added), nil
}
