package config

import "fmt"

const (
	telegramAPIKey = "TELEGRAM_API_KEY"
	redisHost      = "REDIS_HOST"
	redisPort      = "REDIS_PORT"
	redisPassword  = "REDIS_PASSWORD"
	serverHost     = "QUIZ_SERVER_HOST"
	serverPort     = "QUIZ_SERVER_PORT"
	serverLogin    = "QUIZ_SERVER_LOGIN"
	serverPassword = "QUIZ_SERVER_PASSWORD"
)

type Config struct {
	TelegramAPIKey string
	Redis          Redis
	Server         Server
}

type Server struct {
	Host     string
	Port     string
	Login    string
	Password string
}

type Redis struct {
	Host     string
	Port     string
	Password string
}

func (cfg Config) GetServerConnectionString() string {
	return fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
}
