package main

import (
	"context"
	"fmt"
	"github.com/dimayasha7123/quiz_service/config"
	"github.com/dimayasha7123/quiz_service/internal/app"
	"github.com/dimayasha7123/quiz_service/internal/db"
	"github.com/dimayasha7123/quiz_service/internal/mw"
	quizApi "github.com/dimayasha7123/quiz_service/internal/quiz_party_api_client"
	"github.com/dimayasha7123/quiz_service/internal/repository"
	pb "github.com/dimayasha7123/quiz_service/pkg/api"
	"github.com/dimayasha7123/quiz_service/utils"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"net/http"
	"os"
)

const (
	configPath = "./config/config.yaml"
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
		utils.Logger.Fatalf("Can't register service from endpoint: %v", err)
	}
	if err := http.ListenAndServe(fmt.Sprintf(":%s", socket.HTTPPort), mux); err != nil {
		utils.Logger.Fatalf("Error while HTTP-server working: %v")
	}

}

func main() {
	utils.InitializeLogger()
	defer utils.SyncLogger()

	b, err := os.ReadFile(configPath)
	if err != nil {
		utils.Logger.Fatalf("Can't read config file: %v", err)
	}

	cfg, err := config.ParseConfig(b)
	if err != nil {
		utils.Logger.Fatalf("Can't parse config: %v", err)
	}

	utils.Logger.Info("Config unmarshalled")
	ctx := context.Background()

	adp, err := db.New(ctx, cfg.Dsn)
	if err != nil {
		utils.Logger.Fatalf("Can't create DB adapter: %v", err)
	}

	utils.Logger.Info("Get DB adapter")

	qserver := app.New(repository.New(adp), quizApi.New(cfg.QuizAPIKey))
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.Socket.Host, cfg.Socket.GrpcPort))
	if err != nil {
		utils.Logger.Fatalw("Failed to listen TCP",
			"err", err,
			"host", cfg.Socket.Host,
			"gRPCport", cfg.Socket.GrpcPort,
		)
	}
	utils.Logger.Infof("Listening TCP at %s:%s", cfg.Socket.Host, cfg.Socket.GrpcPort)

	var opts []grpc.ServerOption
	opts = []grpc.ServerOption{
		grpc.UnaryInterceptor(mw.LogInterceptor),
	}
	utils.Logger.Infof("Create server options")

	grpcServer := grpc.NewServer(opts...)
	utils.Logger.Info("Create gRPC server")

	pb.RegisterQuizServiceServer(grpcServer, qserver)
	utils.Logger.Info("Register gRPC server")

	go func() {
		err = grpcServer.Serve(lis)
		if err != nil {
			utils.Logger.Fatalf("Error while server working: %v", err)
		}
	}()
	utils.Logger.Info("Server running!")

	go runRest(cfg.Socket)
	utils.Logger.Info("HTTP-proxy server running!")

	for true {
	}
}
