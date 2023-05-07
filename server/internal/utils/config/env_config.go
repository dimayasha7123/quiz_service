package config

import (
	"fmt"
	"os"
	"reflect"
)

type envConfig struct {
	environment environment
	env         map[string]string
}

type environment struct {
	SocketHost       string `env:"SOCKET_HOST"`
	SocketGrpcPort   string `env:"SOCKET_GRPC_PORT"`
	SocketHTTPPort   string `env:"SOCKET_HTTP_PORT"`
	PostgresHost     string `env:"POSTGRES_HOST"`
	PostgresPort     string `env:"POSTGRES_PORT"`
	PostgresUser     string `env:"POSTGRES_USER"`
	PostgresPassword string `env:"POSTGRES_PASSWORD"`
	PostgresDBName   string `env:"POSTGRES_DB"`
}

func New(env map[string]string) (*envConfig, error) {
	err := checkEnvMap(env)
	if err != nil {
		return &envConfig{}, err
	}
	return &envConfig{
		environment: environment{},
		env:         env,
	}, nil
}

func getEnvNameList() []string {
	val := reflect.ValueOf(environment{})
	t := val.Type()
	envNameList := make([]string, t.NumField())
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
		_, okMap := env[envName]
		_, okEnv := os.LookupEnv(envName)
		if !okMap && !okEnv {
			noEnv = append(noEnv, envName)
		}
	}
	if len(noEnv) != 0 {
		return fmt.Errorf("Can't find these ENV variables: %+v", noEnv)
	}
	return nil
}

func (ec *envConfig) Get() *Config {
	cf := &Config{}

	envs := environment{}
	envVal := reflect.ValueOf(&envs).Elem()
	envType := envVal.Type()
	for i := 0; i < envType.NumField(); i++ {
		f := envType.Field(i)
		tag := f.Tag.Get("env")
		var fieldValue string
		osEnv, ok := os.LookupEnv(tag)
		if ok {
			fieldValue = osEnv
		} else {
			fieldValue = ec.env[tag]
		}
		envVal.Field(i).SetString(fieldValue)
	}

	cf.Socket.Host = envs.SocketHost
	cf.Socket.HTTPPort = envs.SocketHTTPPort
	cf.Socket.GrpcPort = envs.SocketGrpcPort
	cf.Dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		envs.PostgresHost,
		envs.PostgresPort,
		envs.PostgresUser,
		envs.PostgresPassword,
		envs.PostgresDBName,
	)

	return cf
}
