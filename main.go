package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/raestrada/sappers/config"
	"github.com/google/uuid"
	"github.com/raestrada/sappers/cluster"
	"github.com/raestrada/sappers/members"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var peers []string

// Command line defaults
const (
	DefaultHTTPAddr = ":11000"
	DefaultRaftAddr = ":12000"
)

// Command line parameters
var inmem bool
var httpAddr string
var raftAddr string
var joinAddr string
var nodeID string
var inPeers string
var gossipPort int

func main() {
	cfg := config.GetConfig()

	var logger *zap.Logger
    var err error

    zapLogLevel := zap.ErrorLevel  // Valor por defecto para el nivel de logs
    logLevel := cfg.LogLevel       // Obtener el nivel de logs desde el singleton de configuración

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

    cfgZap := zap.Config{
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
	logger, err =cfgZap.Build()

	if err != nil {
		panic(err)
	}

	logger = logger.With(zap.String("execution-hash", uuid.New().String()))

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
    cfg := config.GetConfig()  // Acceder a la configuración centralizada

	zap.L().Info("Initializing cluster", 
		zap.String("nodeID", cfg.NodeID), 
		zap.Int("gossipPort", cfg.GossipPort), 
		zap.Strings("peers", cfg.Peers),
	)

    var cluster = cluster.Create(members.GossipMemberListFactory{})
    cluster.Init(cfg.Peers)  // Usar peers desde la configuración
}
