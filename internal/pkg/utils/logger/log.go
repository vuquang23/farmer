package logger

import (
	"errors"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"farmer/internal/pkg/config"
	"farmer/internal/pkg/constants"
)

var zLogger *zap.Logger

func InitLogger() error {
	if zLogger != nil {
		return errors.New("logger is already init")
	}

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
	zLogger = logger
	return nil
}

func FromGinCtx(ctx *gin.Context) *zap.Logger {
	logger, exists := ctx.Get(constants.CtxLoggerKey)
	if !exists {
		return zLogger
	}
	return logger.(*zap.Logger)
}

func WithDescription(desc string) *zap.Logger {
	descriptionField := zapcore.Field{
		Key:    constants.CtxDescriptionKey,
		Type:   zapcore.StringType,
		String: desc,
	}
	logger := zLogger.With(descriptionField)
	return logger
}

func BindLoggerToGinNormCtx(ctx *gin.Context, desc string) error {
	descriptionField := zapcore.Field{
		Key:    constants.CtxDescriptionKey,
		Type:   zapcore.StringType,
		String: desc,
	}
	logger := zLogger.With(descriptionField)
	ctx.Set(constants.CtxLoggerKey, logger)
	return nil
}
