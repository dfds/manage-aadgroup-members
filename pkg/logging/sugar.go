package logging

import (
	"go.uber.org/zap"
)

func GetLogger() *zap.SugaredLogger {
	logger, _ := zap.NewDevelopment()
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			return
		}
	}(logger)
	return logger.Sugar()
}
