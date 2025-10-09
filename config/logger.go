package config

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"sync"
)

const (
	fileStd    = "./log/std.log"
	fileError  = "./log/error.log"
	consoleStd = "stdout"
	consoleErr = "stderr"
)

var GlobalLogger *zap.Logger

type Logger interface {
	SetLogger() (*zap.Logger, error)
}

type ProdLogger struct {
	level      string
	outputPath string
}

type DevLogger struct {
	level      string
	outputPath string
}

func NewLogger(level, outputPath string, prod bool) Logger {
	if !prod {
		return &DevLogger{
			level:      level,
			outputPath: outputPath,
		}
	} else {
		return &ProdLogger{
			level:      level,
			outputPath: outputPath,
		}
	}
}

func (logger *ProdLogger) SetLogger() (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	cfg.EncoderConfig.TimeKey = "time"
	cfg.EncoderConfig.LevelKey = "level"
	cfg.EncoderConfig.MessageKey = "msg"
	cfg.EncoderConfig.CallerKey = "caller"
	switch logger.level {
	case "debug":
		cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn":
		cfg.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		cfg.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	default:
		cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}
	switch logger.outputPath {
	case "console":
		cfg.OutputPaths = []string{consoleStd}
		cfg.ErrorOutputPaths = []string{consoleErr}
	case "file":
		cfg.OutputPaths = []string{fileStd}
		cfg.ErrorOutputPaths = []string{fileError}
	case "double":
		cfg.OutputPaths = []string{consoleStd, fileStd}
		cfg.ErrorOutputPaths = []string{consoleErr, fileError}
	}

	l, err := cfg.Build()
	if err != nil {
		return nil, err
	}
	return l, nil
}

func (logger *DevLogger) SetLogger() (*zap.Logger, error) {
	cfg := zap.NewDevelopmentConfig()
	cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")

	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	cfg.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	switch logger.level {
	case "debug":
		cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn":
		cfg.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		cfg.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	default:
		cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	}
	switch logger.outputPath {
	case "console":
		cfg.OutputPaths = []string{consoleStd}
		cfg.ErrorOutputPaths = []string{consoleErr}
	case "file":
		cfg.OutputPaths = []string{fileStd}
		cfg.ErrorOutputPaths = []string{fileError}
	case "double":
		cfg.OutputPaths = []string{consoleStd, fileStd}
		cfg.ErrorOutputPaths = []string{consoleErr, fileError}
	default:
		cfg.OutputPaths = []string{consoleStd}
		cfg.ErrorOutputPaths = []string{consoleErr}
	}

	l, err := cfg.Build()
	if err != nil {
		return nil, err
	}
	return l, nil
}

func init() {
	var (
		once sync.Once
		err  error
	)
	once.Do(func() {
		GlobalLogger, err = NewLogger("debug", "console", false).SetLogger()
		if err != nil {
			panic(err)
		}
	})
}
