package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func initLoggerConfig(mode string) zap.Config {
	if mode == "DEV" {
		cfg := zap.Config{
			Level:       zap.NewAtomicLevelAt(zap.DebugLevel),
			Development: true,
			Encoding:    "console",
			EncoderConfig: zapcore.EncoderConfig{
				TimeKey:      "ts",
				LevelKey:     "level",
				CallerKey:    "caller",
				EncodeLevel:  zapcore.CapitalColorLevelEncoder,
				MessageKey:   "msg",
				EncodeCaller: zapcore.ShortCallerEncoder,
				EncodeTime:   zapcore.ISO8601TimeEncoder,
			},
			OutputPaths:      []string{"stdout"},
			ErrorOutputPaths: []string{"stderr"},
		}

		return cfg
	}
	cfg := zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.DebugLevel),
		Development: false,
		Encoding:    "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:      "ts",
			LevelKey:     "level",
			CallerKey:    "caller",
			EncodeLevel:  zapcore.LowercaseLevelEncoder,
			MessageKey:   "msg",
			EncodeCaller: zapcore.ShortCallerEncoder,
			EncodeTime:   zapcore.ISO8601TimeEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	return cfg
}
