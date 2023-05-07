package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/dimayasha7123/quiz_service/internal/app"
	"github.com/dimayasha7123/quiz_service/internal/db"
	"github.com/dimayasha7123/quiz_service/internal/mw"
	quizApi "github.com/dimayasha7123/quiz_service/internal/quiz_party_api_client"
	"github.com/dimayasha7123/quiz_service/internal/repository"
	"github.com/dimayasha7123/quiz_service/internal/utils/config"
	"github.com/dimayasha7123/quiz_service/internal/utils/logger"
	"github.com/dimayasha7123/quiz_service/internal/utils/network_config"
	pb "github.com/dimayasha7123/quiz_service/pkg/api"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func runRest(socket config.Socket) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := pb.RegisterQuizServiceHandlerFromEndpoint(
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
	err := logger.RegisterLog()
	if err != nil {
		log.Fatalf("Can't register logger: %v", err)
	}

	var envPath string
	flag.StringVar(&envPath, "env", defaultEnvPath, "path to .env")

	var netCfgPath string
	flag.StringVar(&netCfgPath, "net_config", defaultNetCfgPath, "path to network config")

	flag.Parse()
	logger.Log.Infof("Env path: %s", envPath)
	logger.Log.Infof("Network config path: %s", netCfgPath)

	env, err := godotenv.Read(envPath)
	if err != nil {
		logger.Log.Fatalf("Can't read env file: %v", err)
	}
	logger.Log.Infof("Read env. variables from file: %v", env)

	envCfg, err := config.New(env)
	if err != nil {
		logger.Log.Fatalf("Can't get config from env: %v", err)
	}
	cfg := envCfg.Get()
	logger.Log.Infof("Config unmarshalled: %+v", cfg)

	netCfgKeeper, err := network_config.New(netCfgPath)
	if err != nil {
		logger.Log.Fatalf("Can't get network config: %v", err)
	}
	netCfg := netCfgKeeper.Get()
	logger.Log.Infof("Read network config: %+v", netCfg)

	ctx := context.Background()
	adp, err := db.New(ctx, cfg.Dsn)
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
			"gRPCport", cfg.Socket.GrpcPort,
		)
	}
	logger.Log.Infof("Listening TCP at %s:%s", cfg.Socket.Host, cfg.Socket.GrpcPort)

	opts := []grpc.ServerOption{grpc.UnaryInterceptor(mw.LogInterceptor)}
	logger.Log.Infof("Create server options")

	grpcServer := grpc.NewServer(opts...)
	logger.Log.Info("Create gRPC server")

	pb.RegisterQuizServiceServer(grpcServer, qserver)
	logger.Log.Info("Register gRPC server")

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

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	grpcServer.GracefulStop()
	log.Println("Server exited properly")
}
