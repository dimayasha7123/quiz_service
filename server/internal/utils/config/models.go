package config

const (
	socketHost       = "SOCKET_HOST"
	socketGRPCPort   = "SOCKET_GRPC_PORT"
	socketHTTPPort   = "SOCKET_HTTP_PORT"
	postgresHost     = "POSTGRES_HOST"
	postgresPort     = "POSTGRES_PORT"
	postgresUser     = "POSTGRES_USER"
	postgresPassword = "POSTGRES_PASSWORD"
	postgresDB       = "POSTGRES_DB"
)

type Config struct {
	Socket      Socket
	PostgresDSN string
}

type Socket struct {
	Host     string
	GrpcPort string
	HTTPPort string
}
