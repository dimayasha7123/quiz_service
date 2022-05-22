package config

import (
	"gopkg.in/yaml.v2"
)

type Config struct {
	Socket     Socket
	QuizAPIKey string
	Dsn        string
}

type Socket struct {
	Host     string
	GrpcPort string
	HTTPPort string
}

type configFile struct {
	Socket struct {
		Host     string `yaml:"host"`
		GrpcPort string `yaml:"grpcPort"`
		HTTPPort string `yaml:"httpPort"`
	} `yaml:"socket"`
	QuizAPIKey string `yaml:"quizApiKey"`
	Dsn        string `yaml:"dsn"`
}

func ParseConfig(fileBytes []byte) (*Config, error) {
	cf := configFile{}
	err := yaml.Unmarshal(fileBytes, &cf)
	if err != nil {
		return nil, err
	}

	c := Config{}
	c.QuizAPIKey = cf.QuizAPIKey
	c.Dsn = cf.Dsn
	c.Socket = Socket{
		Host:     cf.Socket.Host,
		GrpcPort: cf.Socket.GrpcPort,
		HTTPPort: cf.Socket.HTTPPort,
	}

	return &c, nil
}
