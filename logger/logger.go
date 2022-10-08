package logger

import (
	"strings"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LoggerConfig struct {
	// Format can be json or console
	Format      string `mapstructure:"format" json:"format" yaml:"format"`
	ServiceName string `mapstructure:"service-name" json:"serviceName" yaml:"service-name"`
	// Level, debug, info, warning, error
	Level         string `mapstructure:"level" json:"level" yaml:"level"`
	DisableCaller bool   `mapstructure:"disable-caller" json:"disableCaller" yaml:"disable-caller"`
	Development   bool   `mapstructure:"development" json:"development" yaml:"development"`
}

var (
	_globalMu sync.RWMutex
	logger    *zap.Logger
)

func L() *zap.Logger {
	_globalMu.RLock()
	l := logger
	_globalMu.RUnlock()
	return l
}

func S() *zap.SugaredLogger {
	return L().Sugar()
}

func MustGetLogger(config *LoggerConfig) *zap.Logger {
	logLevel := parseLogLevel(config.Level)
	// set logger level
	atom := zap.NewAtomicLevelAt(logLevel)

	zapConfig := zap.Config{
		Level:            atom,
		Encoding:         config.Format, // console or json
		EncoderConfig:    getEncoderConfig(),
		ErrorOutputPaths: []string{"stderr"},
		DisableCaller:    config.DisableCaller,
		Development:      config.Development,
		OutputPaths:      []string{"stdout"},
	}
	if config.Format == "json" {
		zapConfig.InitialFields = map[string]interface{}{"S": config.ServiceName}
	}
	zapLogger, err := zapConfig.Build()
	if err != nil {
		panic(err)
	}

	logger = zapLogger

	return zapLogger
}

func parseLogLevel(level string) zapcore.Level {
	switch strings.ToLower(level) {
	case "debug":
		return zap.DebugLevel
	case "info":
		return zap.InfoLevel
	case "warning":
		return zap.WarnLevel
	case "error":
		return zap.ErrorLevel
	default:
		return zap.DebugLevel
	}
}

func getEncoderConfig() zapcore.EncoderConfig {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "T",
		LevelKey:       "L",
		NameKey:        "Logger",
		MessageKey:     "Msg",
		StacktraceKey:  "Stacktrace",
		CallerKey:      "Caller",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder, // 小写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,  // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder, // 全路径编码器
	}
	return encoderConfig
}
