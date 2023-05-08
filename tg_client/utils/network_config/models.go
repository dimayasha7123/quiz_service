package network_config

type Config struct {
	TelegramAPIKey  string `yaml:"telegramApiKey"`
	QuizServiceAuth struct {
		Login    string `yaml:"login"`
		Password int    `yaml:"password"`
	} `yaml:"quiz_service_auth"`
}
