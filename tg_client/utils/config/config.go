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
		Client: Client{
			Host: k.envs[clientHost],
			Port: k.envs[clientPort],
		},
		Redis: Redis{
			Host:     k.envs[redisHost],
			Port:     k.envs[redisPort],
			Password: k.envs[redisPassword],
		},
	}

	return cfg, nil
}

func (k *envConfigKeeper) checkEnvs() error {
	needEnvs := []string{
		clientHost,
		clientPort,
		redisHost,
		redisPort,
		redisPassword,
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
