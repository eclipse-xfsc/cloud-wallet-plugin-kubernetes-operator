package logger

import "go.uber.org/zap"

var logger *zap.Logger

func init() {
	// Creating Logger
	logger = zap.Must(zap.NewDevelopment())
}

func GetLogger() *zap.Logger {
	return logger
}
