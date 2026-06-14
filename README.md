# go_logger 📘

A simple and flexible logger for Golang projects.

## Installation 🚀

```sh
go get github.com/alirezadesh/go_logger
```

## Usage 🛠️

```go
package main

import (
	"github.com/alirezadesh/go_logger"
)

func main() {
  logger := go_logger.New(go_logger.Config{
	  LogToFile:    true,
	  LogToConsole: true,
	  FilePath:     "./logs/",
	  FileName:     "app.log",
	  LogLevel:     go_logger.DebugLevel,
	  ConsoleEncoder: go_logger.EncodeConfig{
		  Time:   true,
		  Caller: true,
		  Level:  true,
	  },
	  FileEncoder: go_logger.EncodeConfig{
		  Time:   true,
		  Caller: true,
		  Level:  true,
	  },
  })

  logger.Info("Logger initialized!")
  logger.Error("something failed", go_logger.Error(err))
}
```

Console output shows colored log levels; the log file is written as plain
text without ANSI color codes.

## Log Levels 📊
 - `DebugLevel`
 - `InfoLevel`
 - `WarnLevel`
 - `ErrorLevel`
 - `DPanicLevel`
 - `PanicLevel`
 - `FatalLevel`

## License 📝
GPL v3

---

Enjoy logging! 🚀

