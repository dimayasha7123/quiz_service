package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/dimayasha7123/quiz_service/config"
	"github.com/dimayasha7123/quiz_service/internal/app"
	"github.com/dimayasha7123/quiz_service/internal/db"
	"github.com/dimayasha7123/quiz_service/internal/mw"
	quizApi "github.com/dimayasha7123/quiz_service/internal/quiz_party_api_client"
	"github.com/dimayasha7123/quiz_service/internal/repository"
	pb "github.com/dimayasha7123/quiz_service/pkg/api"
	"github.com/dimayasha7123/quiz_service/utils/logger"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// TODO: fix this shit
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
		logger.Log.Fatalf("Can't register service from endpoint: %v", err)
	}
	if err := http.ListenAndServe(fmt.Sprintf(":%s", socket.HTTPPort), mux); err != nil {
		logger.Log.Fatalf("Error while HTTP-server working: %v")
	}

}

//docker run --name testPostgres -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=quiz_service_db -d postgres

func main() {
	err := logger.RegisterLog()
	if err != nil {
		log.Fatalf("Can't register logger: %v", err)
	}

	b, err := os.ReadFile(configPath)
	if err != nil {
		logger.Log.Fatalf("Can't read config file: %v", err)
	}

	cfg, err := config.ParseConfig(b)
	if err != nil {
		logger.Log.Fatalf("Can't parse config: %v", err)
	}

	logger.Log.Info("Config unmarshalled")
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

	go func() {
		err = grpcServer.Serve(lis)
		if err != nil {
			logger.Log.Fatalf("Error while server working: %v", err)
		}
	}()
	logger.Log.Info("Server running!")

	go runRest(cfg.Socket)
	logger.Log.Info("HTTP-proxy server running!")

	for true {
	}
}

// что доделать?
// * добить норм логгирование (подумать как его засунуть в проект, пока тупо глобальная переменная)
// * переделать тяжелые запросы в базку
// * переделать конфиг (давать его аргементом мейну)
// * завернуть в докер
// * списков квизов вынести в конфиг??? или не конфиг... сложно...
// *

// что хотелось бы иметь?
// * часть перенести в клик? хотя возможно имеет смысл разбить проект на два сервиса, один из них сохраняет статистику
// * впихнуть метрики
