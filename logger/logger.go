package logger

import (
	"context"
	"github.com/micro/go-micro/v2/metadata"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

type (
	// Logger data
	Logger struct {
		hostname string
		env      string
		zap      *zap.Logger
		ctx      context.Context
		request  interface{}
		response interface{}
		metadata metadata.Metadata
		timing   int64
		sugar    *zap.SugaredLogger
		level    zapcore.Level
	}

	// ILogger interface
	ILogger interface {
		Debug(mgs string)
		Info(mgs string)
		Warn(mgs string)
		Error(mgs string)
		Fatal(mgs string)
		Panic(mgs string)

		Debugf(format string, args ...interface{})
		Infof(format string, args ...interface{})
		Warnf(format string, args ...interface{})
		Errorf(format string, args ...interface{})
		Fatalf(format string, args ...interface{})
		Panicf(format string, args ...interface{})

		Debugw(mgs string, keysAndValues ...interface{})
		Infow(mgs string, keysAndValues ...interface{})
		Warnw(mgs string, keysAndValues ...interface{})
		Errorw(mgs string, keysAndValues ...interface{})
		Fatalw(mgs string, keysAndValues ...interface{})
		Panicw(mgs string, keysAndValues ...interface{})
	}
)

func NewLogger(hostname string, env string) (*Logger, error) {
	var zapLogger *zap.Logger
	var err error
	if env == "production" || env == "LIVE" {
		encoderCfg := zap.NewProductionEncoderConfig()
		encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
		zapLogger = zap.New(zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderCfg),
			zapcore.Lock(os.Stdout),
			zapcore.DebugLevel,
		), zap.AddCaller(), zap.AddCallerSkip(1))
	} else {
		zapLogger, err = zap.NewDevelopment()
	}
	if err != nil {
		return nil, err
	}
	sugar := zapLogger.Sugar()

	return &Logger{
		hostname: hostname,
		env:      env,
		sugar:    sugar,
		zap:      zapLogger,
	}, nil
}
