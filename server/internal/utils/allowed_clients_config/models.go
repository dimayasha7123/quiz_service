package allowed_clients_config

type Config []Client

type Client struct {
	Login    string `yaml:"login"`
	Password string `yaml:"password"`
}
