package config

const (
	quizApiKey       = "QUIZ_API_KEY"
	socketHost       = "SOCKET_HOST"
	socketGRPCPort   = "SOCKET_GRPC_PORT"
	socketHTTPPort   = "SOCKET_HTTP_PORT"
	postgresHost     = "POSTGRES_HOST"
	postgresPort     = "POSTGRES_PORT"
	postgresUser     = "POSTGRES_USER"
	postgresPassword = "POSTGRES_PASSWORD"
	postgresDB       = "POSTGRES_DB"
)

// TODO: можно сделать автоматически через аннотацию типов os.ExpandEnv
type Config struct {
	Socket      Socket
	QuizAPIKey  string
	PostgresDSN string
}

type Socket struct {
	Host     string
	GrpcPort string
	HTTPPort string
}
