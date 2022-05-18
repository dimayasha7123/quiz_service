package repository

import (
	"context"
	"gitlab.ozon.dev/dimayasha7123/homework-2-dimayasha-7123/internal/models"
)

func (r *repository) GetQuizList(ctx context.Context) ([]models.Quiz, error) {

	query := `
		select * from quiz;
	`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	var quizes []models.Quiz

	for rows.Next() {
		var q models.Quiz
		err = rows.Scan(&q.ID, &q.Title)
		if err != nil {
			return nil, err
		}
		quizes = append(quizes, q)
	}

	return quizes, nil
}
