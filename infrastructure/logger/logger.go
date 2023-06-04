package logger

import (
	"fmt"
	"log"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

type ILogger interface {
	Infow(string, ...interface{})
	Errorw(string, ...interface{})
	Warnw(string, ...interface{})
	Debugw(string, ...interface{})
	Fatalw(msg string, keysAndValues ...interface{})
}

func Init(appMode string) *zap.SugaredLogger {
	atom := zap.NewAtomicLevel()

	newLevel, err := getLogLevel(appMode)
	if err == nil {
		atom.SetLevel(newLevel)
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:       "timestamp",
		LevelKey:      "level",
		NameKey:       "logger",
		CallerKey:     "caller",
		FunctionKey:   zapcore.OmitKey,
		MessageKey:    "message",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.LowercaseLevelEncoder,
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.UTC().Format(time.RFC3339Nano))
		},
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	config := zap.Config{
		Level:            atom,
		Encoding:         "json",
		EncoderConfig:    encoderConfig,
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}

	core, err := config.Build()
	if err != nil {
		log.Fatalf("building logger config:%v", err)
	}

	logger := zap.New(core.Core())

	defer logger.Sync()

	sugar := logger.Sugar()

	// setting default fields
	sugar = sugar.With(
		"microservice", "eth_fetcher",
	)

	sugar.Debug("logger initialized")
	return sugar
}

func getLogLevel(appMode string) (zapcore.Level, error) {
	appMode = strings.ToUpper(appMode)
	var newLevel zapcore.Level

	switch appMode {
	case "ERROR":
		return zap.ErrorLevel, nil
	case "WARN":
		return zap.WarnLevel, nil
	case "INFO":
		return zap.InfoLevel, nil
	case "DEBUG":
		return zap.DebugLevel, nil
	default:
		return newLevel, fmt.Errorf("invalid log level %s", appMode)
	}
}
