package logging

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func getLogLevel() zap.AtomicLevel {
	levelStr := os.Getenv("LOG_LEVEL")
	var level zapcore.Level
	if err := level.Set(levelStr); err != nil {
		return zap.NewAtomicLevelAt(zapcore.InfoLevel)
	}
	return zap.NewAtomicLevelAt(level)
}

func CreateLogger() *zap.SugaredLogger {
	config := zap.Config{
		Level:            getLogLevel(),
		Development:      true,
		Encoding:         "console",
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "",
			MessageKey:     "msg",
			StacktraceKey:  "",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalColorLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
	}
	logger, err := config.Build()
	if err != nil {
		panic(err)
	}
	sugar := logger.Sugar()
	return sugar
}
