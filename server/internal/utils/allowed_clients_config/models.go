package allowed_clients_config

type Config []Client

type Client struct {
	Login    string `yaml:"login"`
	Password int    `yaml:"password"`
}
