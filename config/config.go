package config

import (
	"fmt"
	"reflect"

	"gopkg.in/yaml.v2"
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

type configFile struct {
	Socket struct {
		Host     string `yaml:"host" env:"SOCKET_HOST"`
		GrpcPort string `yaml:"grpc_port" env:"SOCKET_GRPC_PORT"`
		HTTPPort string `yaml:"http_port" env:"SOCKET_HTTP_PORT"`
	} `yaml:"socket"`
	QuizAPIKey string `yaml:"quiz_api_key" env:"QUIZ_API_KEY"`
	Postgres   struct {
		Host     string `yaml:"host" env:"POSTGRES_HOST"`
		Port     string `yaml:"port" env:"POSTGRES_PORT"`
		User     string `yaml:"user" env:"POSTGRES_USER"`
		Password string `yaml:"password" env:"POSTGRES_PASSWORD"`
		Dbname   string `yaml:"dbname" env:"POSTGRES_DBNAME"`
		Sslmode  string `yaml:"sslmode" env:"POSTGRES_SSLMODE"`
	} `yaml:"postgres"`
}

func parseConfigBytes(fileBytes []byte) (*configFile, error) {
	cf := configFile{}
	err := yaml.Unmarshal(fileBytes, &cf)
	if err != nil {
		return nil, err
	}
	return &cf, nil
}

func mergeEnvAndConfig(env map[string]string, pval *reflect.Value, pt *reflect.Type) error {
	// нужно пройтись по структуре и посмотреть по тегам, можно ли их вытащить из мапы
	// не будем городить рекурсию чтобы обойти все вложенные структуры, а пока просто если видим,
	// что структура содержит структуру, то проходимся по ней. То есть один уровень вложенности
	val := *pval
	t := *pt
	k := t.Kind()
	fmt.Printf("value: %v\ntype: %v\nkind: %v\n", val, t, k)
	fmt.Printf("NumField: %v\n", t.NumField())
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		fmt.Printf("Kind of field: %v\n", f.Type.Kind())
		if f.Type.Kind() == reflect.Struct {
			nval := reflect.Indirect(reflect.ValueOf(f))
			err := mergeEnvAndConfig(env, &nval, &f.Type)
			if err != nil {
				return err
			}
		} else {
			tag := f.Tag.Get("env")
			envValue, ok := env[tag]
			if ok {
				val.Field(i).SetString(envValue)
				fmt.Printf("Set value <%v> to property <%v>\n", envValue, f.Name)
			}
		}
	}
	return nil
}

func parseConfigFile(cf *configFile) (*Config, error) {
	return nil, nil
}

func GetConfig(configFileBytes []byte, env map[string]string) (*Config, error) {
	cff, err := parseConfigBytes(configFileBytes)
	if err != nil {
		return nil, fmt.Errorf("Can't parse config file bytes: %v", err)
	}

	val := reflect.Indirect(reflect.ValueOf(cff))
	t := val.Type()
	err = mergeEnvAndConfig(env, &val, &t)

	if err != nil {
		return nil, fmt.Errorf("Can't merge config and env: %v", err)
	}
	cf, err := parseConfigFile(cff)
	if err != nil {
		return nil, fmt.Errorf("Can't parse config file: %v", err)
	}
	// сделать проверки конфига...
	return cf, nil
}
