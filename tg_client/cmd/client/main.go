package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	pb "github.com/dimayasha7123/quiz_service/server/pkg/api"
	"github.com/dimayasha7123/quiz_service/tg_client/internal/app"
	"github.com/dimayasha7123/quiz_service/tg_client/internal/db"
	"github.com/dimayasha7123/quiz_service/tg_client/internal/repository"
	"github.com/dimayasha7123/quiz_service/tg_client/utils/config"
	"github.com/dimayasha7123/quiz_service/utils/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var (
	hostname = "localhost"
	crtPath  = "certs/client.crt"
	keyPath  = "certs/client.key"
	caPath   = "certs/ca.crt"
	withMTLS = false
)

func main() {
	zapLogger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("Can't create logger: %v", err)
	}
	sugarLogger := zapLogger.Sugar()
	logger.SetLogger(sugarLogger)

	flag.StringVar(&crtPath, "crt_path", crtPath, "path to client.crt file")
	flag.StringVar(&keyPath, "key_path", keyPath, "path to client.key file")
	flag.StringVar(&caPath, "ca_path", caPath, "path to ca.crt file")
	flag.BoolVar(&withMTLS, "with_mTLS", withMTLS, "enable mTLS")
	flag.Parse()

	configKeeper := config.New()
	cfg, err := configKeeper.Get()
	if err != nil {
		logger.Log.Fatalf("Can't get config from env. vars: %v", err)
	}
	logger.Log.Infof("Config unmarshalled: %+v", cfg)

	if withMTLS {
		logger.Log.Infof("mTLS enabled")
	} else {
		logger.Log.Infof("Insecure mode")
	}

	opts := []grpc.DialOption{grpc.WithInsecure()}

	if withMTLS {
		cert, err := tls.LoadX509KeyPair(crtPath, keyPath)
		if err != nil {
			logger.Log.Fatalf("Can't read cert key pair: %v", err)
		}

		certPool := x509.NewCertPool()
		ca, err := os.ReadFile(caPath)
		if err != nil {
			logger.Log.Fatalf("Can't read ca cert: %v", err)
		}

		ok := certPool.AppendCertsFromPEM(ca)
		if !ok {
			logger.Log.Fatalf("Can't append ca cert to pool")
		}
		logger.Log.Infof("Read TLS certs")

		opts = []grpc.DialOption{
			grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
				ServerName:   hostname,
				Certificates: []tls.Certificate{cert},
				RootCAs:      certPool,
			})),
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	adp := db.New(cfg.Redis)
	logger.Log.Info("Get DB adapter")

	conn, err := grpc.Dial(cfg.GetServerConnectionString(), opts...)
	if err != nil {
		logger.Log.Fatalf("Can't create gRPC connection to quiz server: %v", err)
	}
	defer conn.Close()
	logger.Log.Info("Create gRPC client connection")

	client := app.New(repository.New(adp), cfg.TelegramAPIKey, pb.NewQuizServiceClient(conn))
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
