package config

import "fmt"

type envConfigKeeper struct {
	envs map[string]string
}

func New(envs map[string]string) *envConfigKeeper {
	return &envConfigKeeper{envs: envs}
}

func (k *envConfigKeeper) Get() (Config, error) {
	err := k.checkEnvs()
	if err != nil {
		return Config{}, err
	}

	cfg := Config{
		Socket: Socket{
			Host:     k.envs[socketHost],
			GrpcPort: k.envs[socketGRPCPort],
			HTTPPort: k.envs[socketHTTPPort],
		},
		PostgresDSN: fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			k.envs[postgresHost],
			k.envs[postgresPort],
			k.envs[postgresUser],
			k.envs[postgresPassword],
			k.envs[postgresDB],
		),
	}

	return cfg, nil
}

func (k *envConfigKeeper) checkEnvs() error {
	needEnvs := []string{
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
		_, ok := k.envs[env]
		if !ok {
			notFound = append(notFound, env)
		}
	}

	if len(notFound) == 0 {
		return nil
	}

	return fmt.Errorf("can't found these envs: %v", notFound)
}
