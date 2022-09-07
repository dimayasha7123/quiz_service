package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/dimayasha7123/quiz_service/internal/app"
	"github.com/dimayasha7123/quiz_service/internal/db"
	"github.com/dimayasha7123/quiz_service/internal/mw"
	quizApi "github.com/dimayasha7123/quiz_service/internal/quiz_party_api_client"
	"github.com/dimayasha7123/quiz_service/internal/repository"
	"github.com/dimayasha7123/quiz_service/internal/utils/config"
	"github.com/dimayasha7123/quiz_service/internal/utils/logger"
	pb "github.com/dimayasha7123/quiz_service/pkg/api"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

// Rem protoc --go_out=pkg --go_opt=paths=source_relative --go-grpc_out=pkg --go-grpc_opt=paths=source_relative api/api.proto
// Rem protoc --go_out=. --go-grpc_out=. --grpc-gateway_out=. --grpc-gateway_opt generate_unbound_methods=true --openapiv2_out . api.proto
// protoc -I ./api --go_out ./pkg/api --go_opt paths=source_relative --go-grpc_out ./pkg/api --go-grpc_opt paths=source_relative --grpc-gateway_out ./pkg/api --grpc-gateway_opt paths=source_relative --grpc-gateway_opt generate_unbound_methods=true --openapiv2_out ./pkg/api api/api.proto

const (
	defaultEnvPath = "./.env"
)

func main() {
	err := logger.RegisterLog()
	if err != nil {
		log.Fatalf("Can't register logger: %v", err)
	}

	var envPath string
	flag.StringVar(&envPath, "env", defaultEnvPath, "path to .env")

	flag.Parse()
	if envPath == "" {
		logger.Log.Fatal("Env path can't be empty")
	}

	env, err := godotenv.Read(defaultEnvPath)
	if err != nil {
		logger.Log.Fatalf("Can't read env file: %v", err)
	}

	envCfg, err := config.New(env)
	if err != nil {
		logger.Log.Fatalf("Can't get config from env: %v", err)
	}
	cfg := envCfg.Get()

	logger.Log.Infow("Config unmarshalled", "config", cfg)

	ctx := context.Background()
	adp, err := db.New(ctx, cfg.Dsn)
	if err != nil {
		logger.Log.Fatalf("Can't create DB adapter: %v", err)
	}

	logger.Log.Info("Get DB adapter")

	qserver := app.New(repository.New(adp), quizApi.New(cfg.QuizAPIKey))
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

	go runRest(cfg.Socket)
	logger.Log.Info("HTTP-proxy server running!")

	logger.Log.Info("Server running!")
	err = grpcServer.Serve(lis)
	if err != nil {
		logger.Log.Fatalf("Error while server working: %v", err)
	}
}
