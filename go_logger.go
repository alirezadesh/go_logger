// Package go_logger is a small, configurable wrapper around uber-go/zap.
// It builds a *zap.Logger that can write to the console (with colored levels)
// and/or to a file, with independent field control for each destination.
package go_logger

import (
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Log levels, re-exported from zapcore so callers need not import it directly.
const (
	DebugLevel  = zapcore.DebugLevel
	InfoLevel   = zapcore.InfoLevel
	WarnLevel   = zapcore.WarnLevel
	ErrorLevel  = zapcore.ErrorLevel
	DPanicLevel = zapcore.DPanicLevel
	PanicLevel  = zapcore.PanicLevel
	FatalLevel  = zapcore.FatalLevel

	// InvalidLevel sits just above the highest supported level.
	InvalidLevel = FatalLevel + 1
)

const timeLayout = "2006/01/02 15:04:05"

// Config controls where logs go and how each destination is formatted.
type Config struct {
	FilePath       string
	FileName       string
	LogLevel       zapcore.Level
	LogToFile      bool
	LogToConsole   bool
	ConsoleEncoder EncodeConfig
	FileEncoder    EncodeConfig
}

// EncodeConfig toggles which fields appear in a log entry.
type EncodeConfig struct {
	Time   bool
	Level  bool
	Caller bool
}

// New builds a logger from config. It panics if a requested log file or its
// directory cannot be created.
func New(config Config) *zap.Logger {
	var cores []zapcore.Core

	if config.LogToConsole {
		cores = append(cores, newConsoleCore(config))
	}
	if config.LogToFile {
		cores = append(cores, newFileCore(config))
	}

	return zap.New(zapcore.NewTee(cores...), zap.AddCaller())
}

// Error wraps an error as a structured "error" field.
func Error(err error) zap.Field {
	return zap.NamedError("error", err)
}

// newConsoleCore writes to stdout with colored level names.
func newConsoleCore(config Config) zapcore.Core {
	encoderConfig := newEncoderConfig(config.ConsoleEncoder, colorLevelEncoder)
	encoder := zapcore.NewConsoleEncoder(encoderConfig)
	return zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), config.LogLevel)
}

// newFileCore appends to config.FilePath/config.FileName. Level names are
// written as plain text (no ANSI color codes) so log files stay readable.
func newFileCore(config Config) zapcore.Core {
	if err := os.MkdirAll(config.FilePath, os.ModePerm); err != nil {
		panic("go_logger: failed to create log directory: " + err.Error())
	}

	path := filepath.Join(config.FilePath, config.FileName)
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic("go_logger: failed to open log file: " + err.Error())
	}

	encoderConfig := newEncoderConfig(config.FileEncoder, zapcore.CapitalLevelEncoder)
	encoder := zapcore.NewConsoleEncoder(encoderConfig)
	return zapcore.NewCore(encoder, zapcore.AddSync(file), config.LogLevel)
}

// newEncoderConfig builds an encoder config, dropping any field disabled in
// fields by clearing its key (an empty key tells zap to omit the field).
func newEncoderConfig(fields EncodeConfig, levelEncoder zapcore.LevelEncoder) zapcore.EncoderConfig {
	cfg := zapcore.EncoderConfig{
		TimeKey:      "timestamp",
		LevelKey:     "level",
		MessageKey:   "msg",
		CallerKey:    "caller",
		EncodeTime:   zapcore.TimeEncoderOfLayout(timeLayout),
		EncodeLevel:  levelEncoder,
		EncodeCaller: zapcore.FullCallerEncoder,
	}

	if !fields.Time {
		cfg.TimeKey = ""
	}
	if !fields.Level {
		cfg.LevelKey = ""
	}
	if !fields.Caller {
		cfg.CallerKey = ""
	}
	return cfg
}

// colorLevelEncoder renders the level name in a level-specific color.
func colorLevelEncoder(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(levelColor(level).Sprint(level.CapitalString()))
}

func levelColor(level zapcore.Level) *color.Color {
	switch level {
	case DebugLevel:
		return color.New(color.FgCyan)
	case InfoLevel:
		return color.New(color.FgGreen)
	case WarnLevel:
		return color.New(color.FgYellow)
	case ErrorLevel:
		return color.New(color.FgRed)
	case DPanicLevel, PanicLevel, FatalLevel:
		return color.New(color.FgMagenta)
	default:
		return color.New(color.FgWhite)
	}
}
