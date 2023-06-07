package mw

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/dimayasha7123/quiz_service/server/internal/utils/allowed_clients_config"
	"github.com/dimayasha7123/quiz_service/utils/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"strings"
)

type authKeeper struct {
	clients map[string]string
}

func NewAuthKeeper(config allowed_clients_config.Config) (authKeeper, error) {
	clients := make(map[string]string, len(config))
	doubledLogins := make([]string, 0)

	for _, client := range config {
		_, ok := clients[client.Login]
		if ok {
			doubledLogins = append(doubledLogins, client.Login)
			continue
		}
		clients[client.Login] = client.Password
	}

	if len(doubledLogins) != 0 {
		return authKeeper{}, fmt.Errorf("find doubled logins: %v", doubledLogins)
	}

	return authKeeper{clients: clients}, nil
}

func (k authKeeper) AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		errMsg := "missing metadata"
		logger.Log.Infow("auth interceptor", "error", errMsg)
		return nil, status.Errorf(codes.InvalidArgument, errMsg)
	}

	authMD := md.Get("authorization")
	if len(authMD) == 0 {
		errMsg := "missing authorization in metadata"
		logger.Log.Infow("auth interceptor", "error", errMsg)
		return nil, status.Errorf(codes.InvalidArgument, errMsg)
	}

	auth, err := k.unmarshallAuthMD(authMD)
	if err != nil {
		errMsg := fmt.Sprintf("can't unmarshall authorization in metadata: %v", err)
		logger.Log.Infow("auth interceptor", "error", errMsg)
		return nil, status.Errorf(codes.InvalidArgument, errMsg)
	}

	validErr := k.valid(auth)

	if validErr != nil {
		errMsg := fmt.Sprintf("invalid credentials: %v", validErr)
		logger.Log.Infow("auth interceptor",
			"error", errMsg,
			"login", auth.username,
			"password", auth.password,
		)
		return nil, status.Errorf(codes.Unauthenticated, errMsg)
	}

	logger.Log.Infow("auth interceptor",
		"error", "no",
		"login", auth.username,
		"password", auth.password,
	)

	return handler(ctx, req)
}

type basicAuth struct {
	username string
	password string
}

func (k authKeeper) unmarshallAuthMD(authMD []string) (basicAuth, error) {
	if len(authMD) < 1 {
		return basicAuth{}, fmt.Errorf("no value")
	}

	token := strings.TrimPrefix(authMD[0], "Basic ")
	decodedTokenBytes, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return basicAuth{}, fmt.Errorf("can't decode value")
	}

	data := strings.Split(string(decodedTokenBytes), ":")
	if len(data) < 2 {
		return basicAuth{}, fmt.Errorf("not enought data in token")
	}

	return basicAuth{
		username: data[0],
		password: data[1],
	}, nil
}

func (k authKeeper) valid(auth basicAuth) error {
	pass, ok := k.clients[auth.username]
	if !ok {
		return fmt.Errorf("no such user")
	}
	if pass != auth.password {
		return fmt.Errorf("wrong password")
	}
	return nil
}
