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
		panic(err)
	}
	if err := http.ListenAndServe(fmt.Sprintf(":%s", socket.HTTPPort), mux); err != nil {
		panic(err)
	}
}

func main() {
	b, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatal(err)
	}

	cfg, err := config.ParseConfig(b)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Config unmarshalled")
	ctx := context.Background()

	adp, err := db.New(ctx, cfg.Dsn)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Get db adapter")

	qserver := app.New(repository.New(adp), quizApi.New(cfg.QuizAPIKey))
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.Socket.Host, cfg.Socket.GrpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("Listening TCP at %s:%s", cfg.Socket.Host, cfg.Socket.GrpcPort)

	var opts []grpc.ServerOption
	opts = []grpc.ServerOption{
		grpc.UnaryInterceptor(mw.LogInterceptor),
	}

	log.Println("Create server options")

	grpcServer := grpc.NewServer(opts...)

	log.Println("Create grpc server")

	pb.RegisterQuizServiceServer(grpcServer, qserver)

	log.Println("Register grpc server")
	log.Println("Server running!")

	go runRest(cfg.Socket)
	log.Println("HTTP-proxy server running!")

	err = grpcServer.Serve(lis)
	if err != nil {
		panic(err)
	}
}
