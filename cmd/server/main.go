package main

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"gitlab.ozon.dev/dimayasha7123/homework-2-dimayasha-7123/config"
	"gitlab.ozon.dev/dimayasha7123/homework-2-dimayasha-7123/internal/app"
	"gitlab.ozon.dev/dimayasha7123/homework-2-dimayasha-7123/internal/db"
	"gitlab.ozon.dev/dimayasha7123/homework-2-dimayasha-7123/internal/mw"
	quizApi "gitlab.ozon.dev/dimayasha7123/homework-2-dimayasha-7123/internal/quiz_party_api_client"
	"gitlab.ozon.dev/dimayasha7123/homework-2-dimayasha-7123/internal/repository"
	pb "gitlab.ozon.dev/dimayasha7123/homework-2-dimayasha-7123/pkg/api"
	"gitlab.ozon.dev/dimayasha7123/homework-2-dimayasha-7123/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
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
		log.Fatal(err)
	}
	if err := http.ListenAndServe(fmt.Sprintf(":%s", socket.HTTPPort), mux); err != nil {
		log.Fatal(err)
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

	go runRest(cfg.Socket)
	utils.Logger.Info("HTTP-proxy server running!")

	go func() {
		err = grpcServer.Serve(lis)
		if err != nil {
			utils.Logger.Fatalf("Error while server working: %v", err)
		}
	}()
	utils.Logger.Info("Server running!")

	for true {
	}
}
