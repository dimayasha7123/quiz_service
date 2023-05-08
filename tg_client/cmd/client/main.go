package main

import (
	"context"
	"flag"
	"fmt"
	pb "github.com/dimayasha7123/quiz_service/server/pkg/api"
	"github.com/dimayasha7123/quiz_service/tg_client/internal/app"
	"github.com/dimayasha7123/quiz_service/tg_client/internal/db"
	"github.com/dimayasha7123/quiz_service/tg_client/internal/repository"
	"github.com/dimayasha7123/quiz_service/tg_client/utils/config"
	"github.com/dimayasha7123/quiz_service/tg_client/utils/network_config"
	"github.com/dimayasha7123/quiz_service/utils/env_mapper"
	"github.com/dimayasha7123/quiz_service/utils/logger"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const (
	defaultEnvPath       = ".env"
	defaultClientCfgPath = "tg_bot_config.yaml"
)

func main() {
	zapLogger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("Can't create logger: %v", err)
	}
	sugarLogger := zapLogger.Sugar()
	logger.SetLogger(sugarLogger)

	var envPath string
	flag.StringVar(&envPath, "env", defaultEnvPath, "path to .env")

	var netCfgPath string
	flag.StringVar(&netCfgPath, "config", defaultClientCfgPath, "path to network config")

	flag.Parse()
	logger.Log.Infof("Env path: %s", envPath)
	logger.Log.Infof("Network config path: %s", netCfgPath)

	// TODO: вынести в функцию получение конфига
	fileEnvs, err := godotenv.Read(envPath)
	if err != nil {
		logger.Log.Fatalf("Can't read env file: %v", err)
	}
	logger.Log.Infof("Read env. variables from file: %v", fileEnvs)

	updatedEnvs, updatedCount := env_mapper.ReplaceFileWithEnv(fileEnvs)
	if updatedCount != 0 {
		logger.Log.Infof("Update %d env. variables: %v", updatedCount, updatedEnvs)
	}

	envConfigKeeper := config.New(updatedEnvs)
	cfg, err := envConfigKeeper.Get()
	if err != nil {
		logger.Log.Fatalf("Can't get config from env: %v", err)
	}
	logger.Log.Infof("Config unmarshalled: %+v", cfg)

	// TODO: вынести в функцию получение нетворк конфига
	newCfgBytes, err := os.ReadFile(netCfgPath)
	if err != nil {
		logger.Log.Fatalf("Can't read client config file: %v", err)
	}
	netCfgKeeper := network_config.New(newCfgBytes)
	netCfg, err := netCfgKeeper.Get()
	if err != nil {
		logger.Log.Fatalf("Can't get client config: %v", err)
	}
	logger.Log.Infof("Read client config: %+v", netCfg)

	ctx, cancel := context.WithCancel(context.Background())
	adp := db.New(cfg.Redis)
	logger.Log.Info("Get DB adapter")

	conn, err := grpc.Dial(cfg.GetClientConnectionString(), grpc.WithInsecure())
	if err != nil {
		logger.Log.Fatalf("Can't create gRPC connection to quiz server: %v", err)
	}
	defer conn.Close()
	logger.Log.Info("Create gRPC client connection")

	client := app.New(repository.New(adp), netCfg.TelegramAPIKey, pb.NewQuizServiceClient(conn))
	logger.Log.Info("Create telegram bot handler for quiz server")

	go func() {
		err = client.Run(ctx)
		if err != nil {
			logger.Log.Fatalf("Error while running bot handler: %v", err)
		}
	}()
	logger.Log.Info("Bot handler is running!")

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-done
	fmt.Println()
	logger.Log.Info("Client has been stopped!")

	cancel()
	logger.Log.Info("Client exited properly")
}
