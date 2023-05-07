package config

type Configuration interface {
	Get() *Config
}

type Config struct {
	Socket Socket
	Dsn    string
}

type Socket struct {
	Host     string
	GrpcPort string
	HTTPPort string
}
