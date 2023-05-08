package logger

import (
	"fmt"
	"go.uber.org/zap"
)

func RegisterLog() error {
	fmt.Println("!!!in register log!!!")
	zLogger, err := initLog()
	if err != nil {
		return err
	}
	defer zLogger.Sync()
	zSugarlog := zLogger.Sugar()
	SetLogger(zSugarlog)
	return nil
}

func initLog() (*zap.Logger, error) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}
	return logger, nil
}
