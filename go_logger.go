package go_logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path/filepath"
)

type Config struct {
	LogToFile   bool
	FilePath    string
	FileName    string
	LogLevel    zapcore.Level
	ShowConsole bool
}

const (
	DebugLevel zapcore.Level = iota - 1
	InfoLevel
	WarnLevel
	ErrorLevel
	DPanicLevel
	PanicLevel
	FatalLevel

	_minLevel = DebugLevel
	_maxLevel = FatalLevel

	InvalidLevel = _maxLevel + 1
)

func customColorLevelEncoder(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	colors := map[zapcore.Level]string{
		zapcore.DebugLevel:  "\033[36m",
		zapcore.InfoLevel:   "\033[32m",
		zapcore.WarnLevel:   "\033[33m",
		zapcore.ErrorLevel:  "\033[31m",
		zapcore.DPanicLevel: "\033[35m",
		zapcore.PanicLevel:  "\033[35m",
		zapcore.FatalLevel:  "\033[35m",
	}
	color, exists := colors[level]
	if !exists {
		color = "\033[0m"
	}
	enc.AppendString(color + level.CapitalString() + "\033[0m")
}

func New(config Config) *zap.Logger {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:      "timestamp",
		LevelKey:     "level",
		MessageKey:   "msg",
		CallerKey:    "caller",
		EncodeTime:   zapcore.TimeEncoderOfLayout("2006/01/02 15:04:05"),
		EncodeLevel:  customColorLevelEncoder,
		EncodeCaller: zapcore.FullCallerEncoder,
	}

	encoder := zapcore.NewConsoleEncoder(encoderConfig)
	var cores []zapcore.Core

	if config.ShowConsole {
		consoleCore := zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), config.LogLevel)
		cores = append(cores, consoleCore)
	}

	if config.LogToFile {
		err := os.MkdirAll(config.FilePath, os.ModePerm)
		if err != nil {
			panic("Failed to create log directory: " + err.Error())
		}

		file, err := os.OpenFile(filepath.Join(config.FilePath, config.FileName), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic("Failed to open log file: " + err.Error())
		}

		fileCore := zapcore.NewCore(encoder, zapcore.AddSync(file), config.LogLevel)
		cores = append(cores, fileCore)
	}

	return zap.New(zapcore.NewTee(cores...), zap.AddCaller())
}

func Error(err error) zap.Field {
	return zap.NamedError("error", err)
}

/* How to use:
logger := go_logger.New(go_logger.Config{
LogToFile:   true,
FilePath:    "./logs/",
FileName:    "app.log",
LogLevel:    go_logger.DebugLevel,
ShowConsole: true,
})

logger.Info("Logger initialized!")
*/
