package repository

import (
	"context"
	"fmt"
	"github.com/dimayasha7123/quiz_service/tg_client_2/internal/domain"
)

func (r repository) GetUserByID(ctx context.Context, id int64) (domain.RepoUser, error) {
	key := keyByID(id)

	user, err := r.getUser(ctx, key)
	if err != nil {
		return domain.RepoUser{}, fmt.Errorf("can't get user by key = %v: %v", key, err)
	}

	return user, nil
}
