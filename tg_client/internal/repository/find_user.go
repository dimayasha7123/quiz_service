package repository

import (
	"context"
	"fmt"
)

func (r repository) FindUser(ctx context.Context, tgID int64) (string, error) {
	key := getKeyByTGID(tgID)
	name, err := r.client.HGet(ctx, key, fieldName).Result()
	if err != nil {
		return "", fmt.Errorf("can't get username by key %d: %v", tgID, err)
	}
	return name, nil
}
