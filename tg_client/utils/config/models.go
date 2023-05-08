package config

import "fmt"

const (
	clientHost    = "TG_CLIENT_HOST"
	clientPort    = "TG_CLIENT_PORT"
	redisHost     = "REDIS_HOST"
	redisPort     = "REDIS_PORT"
	redisPassword = "REDIS_PASSWORD"
)

type Config struct {
	Client Client
	Redis  Redis
}

type Client struct {
	Host string
	Port string
}

type Redis struct {
	Host     string
	Port     string
	Password string
}

func (cfg Config) GetClientConnectionString() string {
	return fmt.Sprintf("%s:%s", cfg.Client.Host, cfg.Client.Port)
}
