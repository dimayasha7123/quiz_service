package repository

import (
	"context"
)

func (r *repository) AddUser(ctx context.Context, name string) (int64, error) {
	query := `
		insert into user_account (
			name
		) values (
			$1
		) returning id;
	`
	var ID int64
	err := r.pool.QueryRow(ctx, query, name).Scan(&ID)
	return ID, err
}
