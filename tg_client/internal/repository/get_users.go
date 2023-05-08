package repository

import (
	"context"
	"fmt"
	"strconv"

	"github.com/dimayasha7123/quiz_service/tg_client/internal/models"
)

const defaultCount = 10

func (r repository) GetUsers(ctx context.Context) ([]models.User, error) {
	users := make([]models.User, 0)
	errs := make([]error, 0)

	var cursor uint64
	firstIter := true
	for cursor != 0 || firstIter {
		firstIter = false
		keys, newCursor, err := r.client.Scan(
			ctx,
			cursor,
			fmt.Sprintf("%s*", keyPrefix),
			defaultCount,
		).Result()

		cursor = newCursor
		if err != nil {
			errs = append(errs, fmt.Errorf("can't get keys from redis: %v", err))
			continue
		}

		for _, key := range keys {
			user, err := r.getUser(ctx, key)
			if err != nil {
				errs = append(errs, fmt.Errorf("can't get user by key = %s: %v", key, err))
				continue
			}
			users = append(users, user)
		}
	}
	return users, nil
}

func (r repository) getUser(ctx context.Context, key string) (models.User, error) {
	values, err := r.client.HGetAll(ctx, key).Result()
	if err != nil {
		return models.User{}, err
	}

	tgid, err := getTGIDByKey(key)
	if err != nil {
		return models.User{}, fmt.Errorf("can't get tgid by key = %s: %v", key, err)
	}

	qsidStr, ok := values[fieldQSID]
	if !ok {
		return models.User{}, fmt.Errorf("can't get qsid for user by key = %s", key)
	}
	qsid, err := strconv.Atoi(qsidStr)
	if err != nil {
		return models.User{}, fmt.Errorf(
			"can't convert qsid from string = %s to int by key = %s: %v",
			qsidStr,
			key,
			err,
		)
	}

	username, ok := values[fieldName]
	if !ok {
		return models.User{}, fmt.Errorf("can't get name for user by key = %s", key)
	}

	user := models.User{
		TGID: tgid,
		QSID: int64(qsid),
		Name: username,
	}
	return user, nil
}
