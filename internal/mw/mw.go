package mw

import (
	"context"

	"github.com/dimayasha7123/quiz_service/internal/utils/logger"
	"google.golang.org/grpc"
)

func LogInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	resp, err := handler(ctx, req)
	if err != nil {
		logger.Log.Infow(
			"log interceptor",
			"full method", info.FullMethod,
			"request", req,
			"error", err,
		)
	} else {
		logger.Log.Infow(
			"log interceptor",
			"full method", info.FullMethod,
			"request", req,
		)
	}
	return resp, err
}
