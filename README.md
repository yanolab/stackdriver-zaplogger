# stackdriver-zaplogger

stackdriver-zaplogger is a zap logger interface for stackdriver.\
![cicleci](https://img.shields.io/circleci/project/github/yanolab/stackdriver-zaplogger.svg?label=circleci&style=popout)
![license](https://img.shields.io/github/license/yanolab/stackdriver-zaplogger.svg?style=popout)
![goversion](https://img.shields.io/badge/Go-1.12-green.svg)

## How to install

```
go get -u github.com/yanolab/stackdriver-zaplogger
```

## How to use

For example, Output to both console and stackdriver:
```
package main

import (
	"log"
	"os"

	"cloud.google.com/go/logging"
	zaplogger "github.com/yanolab/stackdriver-zaplogger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/net/context"
)

func main() {
	ctx := context.Background()

	projectID := os.Getenv("PROJECT_ID")
	cli, err := logging.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer cli.Close()

	root, err := newStackdriverZapLogger(zap.DebugLevel, projectID, cli)
	if err != nil {
		log.Fatal(err)
	}

	root.Info("start logging", zap.String("PackageName", "stackdriver-zaplogger"))

	logger := zaplogger.NewLogger(root, "testlogger")
	logger.Debug("created a new logger", zap.String("ProjectID", projectID))
}

func newStackdriverZapLogger(level zapcore.Level, projectID string, cli *logging.Client) (*zap.Logger, error) {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig = zaplogger.EncoderConfig
	config.Sampling = nil
	config.Level = zap.NewAtomicLevelAt(level)

	return config.Build(zap.WrapCore(func(core zapcore.Core) zapcore.Core {
		return zapcore.NewTee(core, zaplogger.NewCore(cli, config.Level))
	}))
}
```