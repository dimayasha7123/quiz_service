package repository

import (
	"context"
	"fmt"

	"github.com/dimayasha7123/quiz_service/tg_client/internal/models"
)

func (r repository) AddUser(ctx context.Context, user models.User) error {
	key := getKeyByTGID(user.TGID)
	values := make(map[string]any)
	values[fieldName] = user.Name
	values[fieldQSID] = user.QSID
	err := r.client.HSet(ctx, key, values).Err()
	if err != nil {
		return fmt.Errorf("can't add user to repo: %v", err)
	}
	return nil

}
