package config

import (
	"fmt"
	"reflect"
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

type enviroment struct {
	SocketHost       string `env:"SOCKET_HOST"`
	SocketGrpcPort   string `env:"SOCKET_GRPC_PORT"`
	SocketHTTPPort   string `env:"SOCKET_HTTP_PORT"`
	QuizAPIKey       string `env:"QUIZ_API_KEY"`
	PostgresHost     string `env:"POSTGRES_HOST"`
	PostgresPort     string `env:"POSTGRES_PORT"`
	PostgresUser     string `env:"POSTGRES_USER"`
	PostgresPassword string `env:"POSTGRES_PASSWORD"`
	PostgresDBName   string `env:"POSTGRES_DBNAME"`
	PostgresSslMode  string `env:"POSTGRES_SSLMODE"`
}

func getEnvNameList() []string {
	val := reflect.Indirect(reflect.ValueOf(enviroment{}))
	t := val.Type()
	envNameList := make([]string, t.NumField())
	fmt.Println(t.NumField())
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		tag := f.Tag.Get("env")
		envNameList[i] = tag
	}
	return envNameList
}

func checkEnvMap(env map[string]string) error {
	noEnv := make([]string, 0)
	for _, envName := range getEnvNameList() {
		_, ok := env[envName]
		if !ok {
			noEnv = append(noEnv, envName)
		}
	}
	if len(noEnv) != 0 {
		return fmt.Errorf("Can't find these ENV variables: %+v", noEnv)
	}
	return nil
}

func GetConfig(env map[string]string) (*Config, error) {
	cf := &Config{}

	err := checkEnvMap(env)
	if err != nil {
		return cf, err
	}

	envs := enviroment{}
	envVal := reflect.ValueOf(&envs).Elem()
	envType := envVal.Type()
	for i := 0; i < envType.NumField(); i++ {
		f := envType.Field(i)
		tag := f.Tag.Get("env")
		envVal.Field(i).SetString(env[tag])
	}

	cf.Socket.Host = envs.SocketHost
	cf.Socket.HTTPPort = envs.SocketHTTPPort
	cf.Socket.GrpcPort = envs.SocketGrpcPort
	cf.QuizAPIKey = envs.QuizAPIKey
	cf.Dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		envs.PostgresHost,
		envs.PostgresPort,
		envs.PostgresUser,
		envs.PostgresPassword,
		envs.PostgresDBName,
		envs.PostgresSslMode,
	)
	
	return cf, nil
}
