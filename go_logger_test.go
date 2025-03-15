package go_logger

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoggerInitialization(t *testing.T) {
	logger := New(Config{
		LogToFile:   false,
		ShowConsole: false,
		LogLevel:    DebugLevel,
	})

	if logger == nil {
		t.Fatal("Logger instance is nil")
	}
}

func TestLoggerFileCreation(t *testing.T) {
	logDir := "./test_logs"
	logFile := "test.log"

	defer func(path string) {
		err := os.RemoveAll(path)
		if err != nil {
			t.Fatal(err)
		}
	}(logDir)

	logger := New(Config{
		LogToFile:   true,
		FilePath:    logDir,
		FileName:    logFile,
		LogLevel:    DebugLevel,
		ShowConsole: false,
	})

	if logger == nil {
		t.Fatal("Logger instance is nil")
	}

	filePath := filepath.Join(logDir, logFile)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Fatalf("Log file %s was not created", filePath)
	}
}

func TestLoggerWritesToFile(t *testing.T) {
	logDir := "./test_logs"
	logFile := "test.log"

	defer func(path string) {
		err := os.RemoveAll(path)
		if err != nil {
			t.Fatal(err)
		}
	}(logDir)

	logger := New(Config{
		LogToFile:   true,
		FilePath:    logDir,
		FileName:    logFile,
		LogLevel:    DebugLevel,
		ShowConsole: false,
	})

	logger.Info("Test log entry")

	filePath := filepath.Join(logDir, logFile)
	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	if len(data) == 0 {
		t.Fatal("Log file is empty, expected log entry")
	}
}
