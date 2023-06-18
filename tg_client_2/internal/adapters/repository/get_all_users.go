package repository

import (
	"context"
	"fmt"
	"github.com/dimayasha7123/quiz_service/tg_client_2/internal/domain"
	"strconv"
)

const defaultCount = 10

func (r repository) GetAllUsers(ctx context.Context) ([]domain.RepoUser, error) {
	users := make([]domain.RepoUser, 0)
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

func (r repository) getUser(ctx context.Context, key string) (domain.RepoUser, error) {
	values, err := r.client.HGetAll(ctx, key).Result()
	if err != nil {
		return domain.RepoUser{}, fmt.Errorf("can't get all values by key = %v: %v", key, err)
	}

	id, err := idByKey(key)
	if err != nil {
		return domain.RepoUser{}, fmt.Errorf("can't get id by key = %s: %v", key, err)
	}

	qsidStr, ok := values[fieldQSID]
	if !ok {
		return domain.RepoUser{}, fmt.Errorf("can't get qsid for user")
	}
	qsid, err := strconv.Atoi(qsidStr)
	if err != nil {
		return domain.RepoUser{}, fmt.Errorf("can't convert qsid = %s from string to int: %v", qsidStr, err)
	}

	username, ok := values[fieldName]
	if !ok {
		return domain.RepoUser{}, fmt.Errorf("can't get name for user")
	}

	user := domain.RepoUser{
		ID:   id,
		QSID: int64(qsid),
		Name: username,
	}
	return user, nil
}
