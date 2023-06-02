package config

import (
	"fmt"
	"os"
)

type configKeeper struct{}

func New() *configKeeper {
	return &configKeeper{}
}

func (k *configKeeper) Get() (Config, error) {
	err := k.checkEnvs()
	if err != nil {
		return Config{}, err
	}

	cfg := Config{
		Socket: Socket{
			Host:     os.Getenv(socketHost),
			GrpcPort: os.Getenv(socketGRPCPort),
			HTTPPort: os.Getenv(socketHTTPPort),
		},
		QuizAPIKey: os.Getenv(quizApiKey),
		PostgresDSN: fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			os.Getenv(postgresHost),
			os.Getenv(postgresPort),
			os.Getenv(postgresUser),
			os.Getenv(postgresPassword),
			os.Getenv(postgresDB),
		),
	}

	return cfg, nil
}

func (k *configKeeper) checkEnvs() error {
	needEnvs := []string{
		quizApiKey,
		socketHost,
		socketGRPCPort,
		socketHTTPPort,
		postgresHost,
		postgresPort,
		postgresUser,
		postgresPassword,
		postgresDB,
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
