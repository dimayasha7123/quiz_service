package repository

import (
	"context"
)

func (r *repository) AddUser(ctx context.Context, name string) (ID int64, err error) {
	query := `
		insert into user_account (
			name
		) values (
			$1
		) returning id;
	`
	err = r.pool.QueryRow(ctx, query, name).Scan(&ID)
	return
}
