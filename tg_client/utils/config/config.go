package config

import (
	"fmt"
	"os"
)

type configKeeper struct {
}

func New() *configKeeper {
	return &configKeeper{}
}

func (k *configKeeper) Get() (Config, error) {
	err := k.checkEnvs()
	if err != nil {
		return Config{}, err
	}

	cfg := Config{
		TelegramAPIKey: os.Getenv(telegramAPIKey),
		Server: Server{
			Host:     os.Getenv(serverHost),
			Port:     os.Getenv(serverPort),
			Login:    os.Getenv(serverLogin),
			Password: os.Getenv(serverPassword),
		},
		Redis: Redis{
			Host:     os.Getenv(redisHost),
			Port:     os.Getenv(redisPort),
			Password: os.Getenv(redisPassword),
		},
	}

	return cfg, nil
}

func (k *configKeeper) checkEnvs() error {
	needEnvs := []string{
		telegramAPIKey,
		redisHost,
		redisPort,
		redisPassword,
		serverHost,
		serverPort,
		serverLogin,
		serverPassword,
	}

	notFound := make([]string, 0, len(needEnvs))
	for _, env := range needEnvs {
		_, ok := os.LookupEnv(env)
		if !ok {
			notFound = append(notFound, env)
		}
	}

	if len(notFound) == 0 {
		return nil
	}

	return fmt.Errorf("can't found these envs: %v", notFound)
}
