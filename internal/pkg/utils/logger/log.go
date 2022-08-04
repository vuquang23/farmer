package logger

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"farmer/internal/pkg/config"
	"farmer/internal/pkg/constants"
)

var Logger *zap.Logger

func InitLogger() error {
	commonCfg := config.Instance().Common
	cfg := zap.Config{
		Level:       zap.NewAtomicLevelAt(zapcore.Level(commonCfg.LogLevel)),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		DisableCaller: true,
		Encoding:      "console",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalColorLevelEncoder,
			EncodeTime:     zapcore.RFC3339TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
		},
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}

	logger, err := cfg.Build()
	if err != nil {
		return err
	}
	Logger = logger
	return nil
}

func FromGinCtx(ctx *gin.Context) *zap.Logger {
	logger, _ := ctx.Get(constants.CtxLoggerKey)
	if logger == nil {
		return Logger
	}
	return logger.(*zap.Logger)
}
