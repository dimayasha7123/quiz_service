package network_config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
)

type NetworkConfiguration interface {
	Get() NetworkConfig
}

type NetworkConfig struct {
	QuizAPIKey     string `yaml:"quiz_api_key"`
	AllowedClients []struct {
		Login    string `yaml:"login"`
		Password int    `yaml:"password"`
	} `yaml:"allowed_clients"`
}

type networkConfigKeeper struct {
	cfg NetworkConfig
}

func (k *networkConfigKeeper) Get() NetworkConfig {
	return k.cfg
}

func New(path string) (*networkConfigKeeper, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("can't read file %s: %v", path, err)
	}

	var cfg NetworkConfig
	err = yaml.Unmarshal(bytes, &cfg)
	if err != nil {
		return nil, fmt.Errorf("can't unmarshall network config: %v", err)
	}

	return &networkConfigKeeper{cfg: cfg}, nil
}
