package repository

import (
	"context"
	"gitlab.ozon.dev/dimayasha7123/homework-2-dimayasha-7123/internal/models"
)

func (r *repository) GetQuiz(ctx context.Context, ID int64) (models.Quiz, error) {
	query := `
		select title
		from quiz
		where id = $1;
	`
	var title string
	err := r.pool.QueryRow(ctx, query, ID).Scan(&title)
	if err != nil {
		return models.Quiz{ID: -1}, err
	}

	return models.Quiz{ID: ID, Name: title}, nil
}
