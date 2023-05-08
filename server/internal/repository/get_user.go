package repository

import (
	"context"
	"github.com/jackc/pgx/v4"
)

func (r *repository) GetUser(ctx context.Context, name string) (int64, error) {
	query := `
		select id from user_account where name = $1;
	`
	var ID int64
	err := r.pool.QueryRow(ctx, query, name).Scan(&ID)
	if err == pgx.ErrNoRows {
		return -1, nil
	}

	return ID, err
}
