package repository

import (
	"context"
	"fmt"
	"github.com/dimayasha7123/quiz_service/tg_client_2/internal/domain"
)

func (r repository) AddUser(ctx context.Context, user domain.RepoUser) error {
	key := keyByID(user.ID)
	values := make(map[string]any)
	values[fieldName] = user.Name
	values[fieldQSID] = user.QSID
	err := r.client.HSet(ctx, key, values).Err()
	if err != nil {
		return fmt.Errorf("can't add user to repo: %v", err)
	}
	return nil

}
