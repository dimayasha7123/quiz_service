package db

import (
	"fmt"
	"github.com/dimayasha7123/quiz_service/tg_client_2/utils/config"
	"github.com/redis/go-redis/v9"
)

func New(cfg config.Redis) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       0,
	})
	return client
}
