package go_logger

import (
	"github.com/fatih/color"
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
	var c *color.Color

	switch level {
	case zapcore.DebugLevel:
		c = color.New(color.FgCyan)
	case zapcore.InfoLevel:
		c = color.New(color.FgGreen)
	case zapcore.WarnLevel:
		c = color.New(color.FgYellow)
	case zapcore.ErrorLevel:
		c = color.New(color.FgRed)
	case zapcore.DPanicLevel, zapcore.PanicLevel, zapcore.FatalLevel:
		c = color.New(color.FgMagenta)
	default:
		c = color.New(color.FgWhite)
	}

	coloredLevel := c.Sprint(level.CapitalString())
	enc.AppendString(coloredLevel)
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
