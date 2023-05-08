package network_config

import (
	"fmt"
	"gopkg.in/yaml.v2"
)

type netConfigKeeper struct {
	bytes []byte
}

func New(bytes []byte) *netConfigKeeper {
	return &netConfigKeeper{bytes: bytes}
}

func (k *netConfigKeeper) Get() (Config, error) {
	var cfg Config

	err := yaml.Unmarshal(k.bytes, &cfg)
	if err != nil {
		return Config{}, fmt.Errorf("can't unmarshall network config: %v", err)
	}

	return cfg, nil
}
