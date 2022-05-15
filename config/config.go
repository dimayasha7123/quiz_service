package config

import (
	"gopkg.in/yaml.v2"
)

type ApiKeys struct {
	Telegram string
	Quiz     string
}

type Config struct {
	ApiKeys ApiKeys
}

type configFile struct {
	APIKeys struct {
		Telegram string `yaml:"telegram"`
		Quiz     string `yaml:"quiz"`
	} `yaml:"apiKeys"`
}

func ParseConfig(fileBytes []byte) (*Config, error) {
	cf := configFile{}
	err := yaml.Unmarshal(fileBytes, &cf)
	if err != nil {
		return nil, err
	}

	c := Config{}

	c.ApiKeys.Telegram = cf.APIKeys.Telegram
	c.ApiKeys.Quiz = cf.APIKeys.Quiz

	return &c, nil
}
