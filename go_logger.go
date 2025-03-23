package go_logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path/filepath"
)

type Config struct {
	FilePath       string
	FileName       string
	LogLevel       zapcore.Level
	LogToFile      bool
	LogToConsole   bool
	ConsoleEncoder EncodeConfig
	FileEncoder    EncodeConfig
}

type EncodeConfig struct {
	Time   bool
	Level  bool
	Caller bool
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

func applyEncoderConfig(enc zapcore.EncoderConfig, config EncodeConfig) zapcore.EncoderConfig {
	if !config.Time {
		enc.TimeKey = ""
	}
	if !config.Level {
		enc.LevelKey = ""
	}
	if !config.Caller {
		enc.CallerKey = ""
	}
	return enc
}

func New(config Config) *zap.Logger {
	baseEncoderConfig := zapcore.EncoderConfig{
		TimeKey:      "timestamp",
		LevelKey:     "level",
		MessageKey:   "msg",
		CallerKey:    "caller",
		EncodeTime:   zapcore.TimeEncoderOfLayout("2006/01/02 15:04:05"),
		EncodeLevel:  customColorLevelEncoder,
		EncodeCaller: zapcore.FullCallerEncoder,
	}

	consoleEncoderConfig := applyEncoderConfig(baseEncoderConfig, config.ConsoleEncoder)
	fileEncoderConfig := applyEncoderConfig(baseEncoderConfig, config.FileEncoder)

	var cores []zapcore.Core

	if config.LogToConsole {
		consoleCore := zapcore.NewCore(zapcore.NewConsoleEncoder(consoleEncoderConfig), zapcore.AddSync(os.Stdout), config.LogLevel)
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

		fileCore := zapcore.NewCore(zapcore.NewConsoleEncoder(fileEncoderConfig), zapcore.AddSync(file), config.LogLevel)
		cores = append(cores, fileCore)
	}

	return zap.New(zapcore.NewTee(cores...), zap.AddCaller())
}

func Error(err error) zap.Field {
	return zap.NamedError("error", err)
}

/* How to use:
logger := go_logger.New(Config{
		LogToFile:    true,
		LogToConsole: true,
		FilePath:     "./logs/",
		FileName:     "app.log",
		LogLevel:     DebugLevel,
		ConsoleEncoder: EncodeConfig{
			Time:   true,
			Caller: true,
			Level:  true,
		}, FileEncoder: EncodeConfig{
			Time:   true,
			Caller: true,
			Level:  true,
		},
	})

logger.Info("Logger initialized!")
*/
