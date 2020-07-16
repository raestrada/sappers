package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/google/uuid"
	"github.com/raestrada/sappers/cluster"
	"github.com/raestrada/sappers/member"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	var logger *zap.Logger

	logLevel := "Error"
	if value, ok := os.LookupEnv("LOG_LEVEL"); ok {
		logLevel = value
	}

	var err error
	zapLogLevel := zap.DebugLevel

	switch logLevel {
	case "INFO":
		zapLogLevel = zap.InfoLevel
	case "DEBUG":
		zapLogLevel = zap.DebugLevel
	case "WARN":
		zapLogLevel = zap.WarnLevel
	case "ERROR":
		zapLogLevel = zap.ErrorLevel
	default:
		zapLogLevel = zap.ErrorLevel
	}

	cfg := zap.Config{
		Encoding:         "json",
		Level:            zap.NewAtomicLevelAt(zapLogLevel),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stdout"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "msg",

			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalLevelEncoder,

			TimeKey:    "time",
			EncodeTime: zapcore.ISO8601TimeEncoder,

			CallerKey:    "caller",
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}
	logger, err = cfg.Build()

	if err != nil {
		panic(err)
	}

	logger = logger.With(zap.String("hash", uuid.New()))

	zap.ReplaceGlobals(logger)
	zap.L().Info("STDOUT Global Logger started")

	go startCluster()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	kill := make(chan os.Signal, 1)
	signal.Notify(kill, os.Kill)
	select {
	case <-interrupt:
		fmt.Println("Interrupt heard. Ending main function")
		return
	case <-kill:
		fmt.Println("Kill heard. Ending main function")
		return
	}
}

func startCluster() {
	var cluster = cluster.Cluster.Create(member.GossipMemberListFactory)
}
