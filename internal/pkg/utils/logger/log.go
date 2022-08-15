package logger

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"farmer/internal/pkg/config"
	"farmer/internal/pkg/constants"
)

var Logger *zap.Logger

func InitLogger() error {
	if Logger != nil {
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

func WithDescription(desc string) *zap.Logger {
	descriptionField := zapcore.Field{
		Key:    constants.CtxDescriptionKey,
		Type:   zapcore.StringType,
		String: desc,
	}
	logger := Logger.With(descriptionField)
	return logger
}

func BindLoggerToGinNormCtx(ctx *gin.Context, desc string) error {
	descriptionField := zapcore.Field{
		Key:    constants.CtxDescriptionKey,
		Type:   zapcore.StringType,
		String: desc,
	}
	logger := Logger.With(descriptionField)
	ctx.Set(constants.CtxLoggerKey, logger)
	return nil
}

func BindLoggerToGinReqCtx(c *gin.Context) error {
	requestIDField := zapcore.Field{
		Key:    constants.CtxRequestIDKey,
		Type:   zapcore.StringType,
		String: uuid.New().String(),
	}

	builder := strings.Builder{}
	builder.WriteString(c.Request.Method)
	builder.WriteString(" ")
	builder.WriteString(c.Request.URL.Path)
	raw := c.Request.URL.RawQuery
	if raw != "" {
		builder.WriteString("?")
		builder.WriteString(raw)
	}
	apiField := zapcore.Field{
		Key:    constants.CtxAPIRequestKey,
		Type:   zapcore.StringType,
		String: builder.String(),
	}

	logger := Logger.
		With(requestIDField).
		With(apiField)

	c.Set(constants.CtxLoggerKey, logger)

	return nil
}
