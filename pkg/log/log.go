package log

import (
	"go.uber.org/zap"
)

type Log struct {
	*zap.SugaredLogger
}

type Logger interface {
	Debugf(pattern string, msg ...interface{})
	Infof(pattern string, msg ...interface{})
	Errorf(pattern string, msg ...interface{})
	Fatalf(pattern string, msg ...interface{})
}

func NewLogger() (*Log, error) {
	cfg := initLoggerConfig()
	logger, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	defer logger.Sync()

	return &Log{
		logger.Sugar(),
	}, nil
}
