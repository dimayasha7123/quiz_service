package network_config

type Config struct {
	QuizAPIKey     string `yaml:"quiz_api_key"`
	AllowedClients []struct {
		Login    string `yaml:"login"`
		Password int    `yaml:"password"`
	} `yaml:"allowed_clients"`
}
