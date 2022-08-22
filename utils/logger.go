package utils

import (
	"go.uber.org/zap"
)

var Logger *zap.SugaredLogger

func InitializeLogger() {
	logger, _ := zap.NewDevelopment()
	Logger = logger.Sugar()
}

func SyncLogger() {
	Logger.Sync()
}
