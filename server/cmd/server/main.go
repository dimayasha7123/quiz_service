package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/dimayasha7123/quiz_service/server/internal/app"
	"github.com/dimayasha7123/quiz_service/server/internal/db"
	"github.com/dimayasha7123/quiz_service/server/internal/mw"
	"github.com/dimayasha7123/quiz_service/server/internal/quiz_party_api_client"
	"github.com/dimayasha7123/quiz_service/server/internal/repository"
	"github.com/dimayasha7123/quiz_service/server/internal/utils/config"
	"github.com/dimayasha7123/quiz_service/server/internal/utils/network_config"
	"github.com/dimayasha7123/quiz_service/server/pkg/api"
	"github.com/dimayasha7123/quiz_service/utils/env_mapper"
	"github.com/dimayasha7123/quiz_service/utils/logger"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func runRest(socket config.Socket) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := api.RegisterQuizServiceHandlerFromEndpoint(
		ctx,
		mux,
		fmt.Sprintf("%s:%s", socket.Host, socket.GrpcPort),
		opts,
	)
	if err != nil {
		logger.Log.Fatalf("Can't register service from endpoint: %v", err)
	}
	if err := http.ListenAndServe(fmt.Sprintf(":%s", socket.HTTPPort), mux); err != nil {
		logger.Log.Fatalf("Error while HTTP-server working: %v")
	}
}

const (
	defaultEnvPath    = ".env"
	defaultNetCfgPath = "network_config.yaml"
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
	flag.StringVar(&netCfgPath, "config", defaultNetCfgPath, "path to network config")

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
		logger.Log.Fatalf("Can't read network config file: %v", err)
	}
	netCfgKeeper := network_config.New(newCfgBytes)
	netCfg, err := netCfgKeeper.Get()
	if err != nil {
		logger.Log.Fatalf("Can't get network config: %v", err)
	}
	logger.Log.Infof("Read network config: %+v", netCfg)

	ctx := context.Background()
	adp, err := db.New(ctx, cfg.PostgresDSN)
	if err != nil {
		logger.Log.Fatalf("Can't create DB adapter: %v", err)
	}
	logger.Log.Info("Get DB adapter")

	qserver := app.New(repository.New(adp), quizApi.New(netCfg.QuizAPIKey))
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.Socket.Host, cfg.Socket.GrpcPort))
	if err != nil {
		logger.Log.Fatalw("Failed to listen TCP",
			"err", err,
			"host", cfg.Socket.Host,
			"gRPCPort", cfg.Socket.GrpcPort,
		)
	}
	logger.Log.Infof("Listening TCP at %s:%s", cfg.Socket.Host, cfg.Socket.GrpcPort)

	opts := []grpc.ServerOption{grpc.UnaryInterceptor(mw.LogInterceptor)}
	grpcServer := grpc.NewServer(opts...)
	api.RegisterQuizServiceServer(grpcServer, qserver)
	logger.Log.Info("Create and register gRPC server")

	go func() {
		err = grpcServer.Serve(lis)
		if err != nil && err != grpc.ErrServerStopped {
			logger.Log.Fatalf("Error while server working: %v", err)
		}
	}()
	logger.Log.Info("Server is running!")

	go runRest(cfg.Socket)
	logger.Log.Info("HTTP-proxy server running!")

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-done
	fmt.Println()
	log.Println("Server has been stopped")

	//ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	//defer cancel()

	grpcServer.GracefulStop()
	log.Println("Server exited properly")
}
