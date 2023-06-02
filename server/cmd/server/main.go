package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"github.com/dimayasha7123/quiz_service/server/internal/app"
	"github.com/dimayasha7123/quiz_service/server/internal/db"
	"github.com/dimayasha7123/quiz_service/server/internal/mw"
	"github.com/dimayasha7123/quiz_service/server/internal/quiz_party_api_client"
	"github.com/dimayasha7123/quiz_service/server/internal/repository"
	"github.com/dimayasha7123/quiz_service/server/internal/utils/allowed_clients_config"
	"github.com/dimayasha7123/quiz_service/server/internal/utils/config"
	"github.com/dimayasha7123/quiz_service/server/pkg/api"
	"github.com/dimayasha7123/quiz_service/utils/logger"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

// TODO: сделать похожим на main, может даже вынести в main
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

var (
	allowedClientsCfgPath = "allowed_clients.yaml"
	crtPath               = "certs/server.crt"
	keyPath               = "certs/server.key"
	caPath                = "certs/ca.crt"
	withMTLS              = false
)

func main() {
	zapLogger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("Can't create logger: %v", err)
	}
	sugarLogger := zapLogger.Sugar()
	logger.SetLogger(sugarLogger)

	flag.StringVar(&allowedClientsCfgPath, "allowed_clients_config_path", allowedClientsCfgPath, "path to allowed clients config")
	flag.StringVar(&crtPath, "crt_path", crtPath, "path to server.crt file")
	flag.StringVar(&keyPath, "key_path", keyPath, "path to server.key file")
	flag.StringVar(&caPath, "ca_path", caPath, "path to ca.crt file")
	flag.BoolVar(&withMTLS, "with_mTLS", withMTLS, "enable mTLS")
	flag.Parse()

	configKeeper := config.New()
	cfg, err := configKeeper.Get()
	if err != nil {
		logger.Log.Fatalf("Can't get config from env vars: %v", err)
	}
	logger.Log.Infof("Config unmarshalled: %+v", cfg)

	cfgBytes, err := os.ReadFile(allowedClientsCfgPath)
	if err != nil {
		logger.Log.Fatalf("Can't read allowed clients config file: %v", err)
	}
	cfgKeeper := allowed_clients_config.New(cfgBytes)
	allowedClientsCfg, err := cfgKeeper.Get()
	if err != nil {
		logger.Log.Fatalf("Can't get allowed clients config: %v", err)
	}
	logger.Log.Infof("Read allowed clients config: %+v", allowedClientsCfg)

	if withMTLS {
		logger.Log.Infof("mTLS enabled")
	} else {
		logger.Log.Infof("Insecure mode")
	}

	ctx := context.Background()

	adp, err := db.New(ctx, cfg.PostgresDSN)
	if err != nil {
		logger.Log.Fatalf("Can't create DB adapter: %v", err)
	}
	logger.Log.Info("Get DB adapter")

	qserver := app.New(repository.New(adp), quizApi.New(cfg.QuizAPIKey))
	netAddr := fmt.Sprintf("%s:%s", cfg.Socket.Host, cfg.Socket.GrpcPort)
	lis, err := net.Listen("tcp", netAddr)
	if err != nil {
		logger.Log.Fatalw("Failed to listen TCP",
			"err", err,
			"netAddr", netAddr,
		)
	}
	logger.Log.Infof("Listening TCP at %s", netAddr)

	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(mw.LogInterceptor),
	}

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

		opts = append(opts, grpc.Creds(credentials.NewTLS(&tls.Config{
			ClientAuth:   tls.RequireAndVerifyClientCert,
			Certificates: []tls.Certificate{cert},
			ClientCAs:    certPool,
		})))
	}

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
