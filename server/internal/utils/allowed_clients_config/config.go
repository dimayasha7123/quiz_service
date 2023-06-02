package allowed_clients_config

import (
	"fmt"
	"gopkg.in/yaml.v2"
)

type configKeeper struct {
	bytes []byte
}

func New(bytes []byte) *configKeeper {
	return &configKeeper{bytes: bytes}
}

func (k *configKeeper) Get() (Config, error) {
	var cfg Config

	err := yaml.Unmarshal(k.bytes, &cfg)
	if err != nil {
		return Config{}, fmt.Errorf("can't unmarshall network config: %v", err)
	}

	return cfg, nil
}
