# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

`go_logger` is a thin wrapper library around [uber-go/zap](https://github.com/uber-go/zap) that returns a configured `*zap.Logger`. The entire public API lives in `go_logger.go`; there are no subpackages.

## Commands

```sh
go test ./...                          # run all tests
go test -run TestLoggerWritesToFile    # run a single test by name
go build ./...                         # compile
go vet ./...                           # static checks
```

Tests in `go_logger_test.go` create and clean up a `./test_logs` directory; they exercise real file I/O rather than mocks.

## Architecture

`New(config Config) *zap.Logger` builds the logger by composing zap cores with `zapcore.NewTee`, delegating to two helpers:

- `newConsoleCore` (when `LogToConsole`) writes to `os.Stdout` using `colorLevelEncoder` (colored level names).
- `newFileCore` (when `LogToFile`) creates `FilePath` via `os.MkdirAll` and appends to `FilePath/FileName`, using `zapcore.CapitalLevelEncoder` so the file stays **plain text with no ANSI color codes**. **Note:** failures to create the directory or open the file `panic` rather than returning an error.

`newEncoderConfig(fields EncodeConfig, levelEncoder)` builds each core's encoder config and takes the level encoder as a parameter — this is what keeps color out of the file while keeping it on the console. It disables a field by blanking the corresponding key (`TimeKey`/`LevelKey`/`CallerKey`); an empty key tells zap to omit the field.

`colorLevelEncoder` / `levelColor` colorize the level string via `fatih/color` (cyan/green/yellow/red/magenta per level), and are used only by the console core.

## Conventions

- Log levels are re-exported as package constants (`DebugLevel`...`FatalLevel`, plus `InvalidLevel`) aliasing `zapcore.Level` values, so callers don't import zapcore directly for levels.
- `Error(err error) zap.Field` is a helper that wraps an error as a `zap.NamedError("error", err)` field for structured logging.
- The module path is `github.com/alirezadesh/go_logger`; the package name is `go_logger`.
