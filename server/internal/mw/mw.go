package mw

import (
	"context"
	"github.com/dimayasha7123/quiz_service/server/internal/utils/logger"

	"google.golang.org/grpc"
)

func LogInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	resp, err := handler(ctx, req)
	msg := "log interceptor"
	if err != nil {
		logger.Log.Infow(
			msg,
			"full method", info.FullMethod,
			"request", req,
			"error", err,
		)
	} else {
		logger.Log.Infow(
			msg,
			"full method", info.FullMethod,
			"request", req,
		)
	}
	return resp, err
}
