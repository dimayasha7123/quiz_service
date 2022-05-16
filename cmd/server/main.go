package main

import (
	"context"
	"gitlab.ozon.dev/dimayasha7123/homework-2-dimayasha-7123/config"
	"gitlab.ozon.dev/dimayasha7123/homework-2-dimayasha-7123/internal/app"
	"gitlab.ozon.dev/dimayasha7123/homework-2-dimayasha-7123/internal/db"
	"gitlab.ozon.dev/dimayasha7123/homework-2-dimayasha-7123/internal/mw"
	"gitlab.ozon.dev/dimayasha7123/homework-2-dimayasha-7123/internal/repository"
	pb "gitlab.ozon.dev/dimayasha7123/homework-2-dimayasha-7123/pkg/api"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

const (
	configPath = "./config/config.yaml"
)

func main() {
	b, errr := os.ReadFile(configPath)
	if errr != nil {
		log.Fatal(errr)
	}

	cfg, errc := config.ParseConfig(b)
	if errc != nil {
		log.Fatal(errc)
	}

	//log.Printf("Config = %+v\n", cfg)
	log.Println("Config unmarshalled")
	ctx := context.Background()

	adp, err := db.New(ctx, cfg.Dsn)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Get db adapter")

	newServer := app.New(repository.New(adp))
	lis, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Println("Listening TCP at localhost:8080")

	var opts []grpc.ServerOption
	opts = []grpc.ServerOption{
		grpc.UnaryInterceptor(mw.LogInterceptor),
	}

	log.Println("Create server options")

	grpcServer := grpc.NewServer(opts...)

	log.Println("Create grpc server")

	pb.RegisterQuizServiceServer(grpcServer, newServer)

	log.Println("Register grpc server")

	log.Println("Server running!")

	err = grpcServer.Serve(lis)
	if err != nil {
		panic(err)
	}
}
