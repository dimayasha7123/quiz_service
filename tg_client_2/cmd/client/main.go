package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"github.com/dimayasha7123/quiz_service/server/pkg/api"
	"github.com/dimayasha7123/quiz_service/tg_client_2/internal/adapters/db"
	"github.com/dimayasha7123/quiz_service/tg_client_2/internal/adapters/repository"
	"github.com/dimayasha7123/quiz_service/tg_client_2/internal/app"
	"github.com/dimayasha7123/quiz_service/tg_client_2/internal/domain"
	"github.com/dimayasha7123/quiz_service/tg_client_2/internal/ports/auth"
	"github.com/dimayasha7123/quiz_service/tg_client_2/internal/ports/telegram_bot"
	"github.com/dimayasha7123/quiz_service/tg_client_2/utils/config"
	"github.com/dimayasha7123/quiz_service/utils/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	hostname      = "localhost"
	crtPath       = "certs/client.crt"
	keyPath       = "certs/client.key"
	caPath        = "certs/ca.crt"
	withMTLS      = false
	withBasicAuth = false
	httpTimeout   = 20 * time.Second
	tgApiDelay    = 50 * time.Millisecond
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
	flag.BoolVar(&withBasicAuth, "with_basic_auth", withBasicAuth, "enable basic auth")

	flag.DurationVar(&httpTimeout, "http_timeout", httpTimeout, "timeout for http client to interact with Telegram Bot API")
	flag.DurationVar(&tgApiDelay, "tg_api_delay", tgApiDelay, "delay for requests to Telegram Bot API")

	flag.Parse()

	configKeeper := config.New()
	cfg, err := configKeeper.Get()
	if err != nil {
		logger.Log.Fatalf("Can't get config from env. vars: %v", err)
	}
	logger.Log.Infof("Config unmarshalled: %+v", cfg)

	opts := []grpc.DialOption{grpc.WithInsecure()}

	if withMTLS {
		logger.Log.Infof("mTLS enabled")

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
	} else {
		logger.Log.Infof("Insecure mode")
	}

	if withBasicAuth {
		if !withMTLS {
			logger.Log.Fatalf("Can't enable basic auth without using mTLS! Enable \"withMTLS\" flag")
		}
		logger.Log.Infof("Basic auth enabled")

		authData := auth.New(cfg.Server.Login, cfg.Server.Password)
		opts = append(opts, grpc.WithPerRPCCredentials(authData))
	} else {
		logger.Log.Infof("Basic auth disabled")
	}

	ctx, cancel := context.WithCancel(context.Background())
	adp := db.New(cfg.Redis)
	logger.Log.Info("Get DB adapter")

	sessions, err := domain.NewSessions(ctx, repository.New(adp))
	if err != nil {
		logger.Log.Fatalf("Can't create users sessions: %v", err)
	}
	logger.Log.Info("Create users sessions")

	conn, err := grpc.Dial(cfg.GetServerConnectionString(), opts...)
	if err != nil {
		logger.Log.Fatalf("Can't create gRPC connection to quiz server: %v", err)
	}
	defer conn.Close()
	logger.Log.Info("Create gRPC client connection")

	app := app.New(sessions, api.NewQuizServiceClient(conn))
	tg_service := telegram_bot.New(
		cfg.TelegramAPIKey,
		app,
		httpTimeout,
		tgApiDelay,
	)
	logger.Log.Info("Create Telegram Bot Service for Quiz Server")

	go func() {
		err = tg_service.Run(ctx)
		if err != nil {
			logger.Log.Fatalf("Error while running bot service: %v", err)
		}
	}()
	logger.Log.Info("Bot service is running!")

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-done
	fmt.Println()
	logger.Log.Info("Bot service has been stopped!")

	cancel()
	logger.Log.Info("Bot service exited properly")

}
