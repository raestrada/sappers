package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"

	"github.com/spf13/viper"
    "github.com/spf13/pflag"
	"github.com/raestrada/sappers/config"
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
	// Definir los parámetros de CLI con pflag
    pflag.Int("gossip-port", 7946, "Puerto para gossip")
    pflag.String("raft-addr", ":12000", "Dirección para Raft")
    pflag.String("http-addr", ":11000", "Dirección HTTP")
    pflag.String("node-id", "default-node", "ID del nodo")
    pflag.String("log-level", "ERROR", "Nivel de logs")
    pflag.StringSlice("peers", []string{"127.0.0.1"}, "Peers del clúster")

    // Parsear los parámetros de CLI
    pflag.Parse()

    // Viper automáticamente tomará los valores de pflag
    viper.BindPFlags(pflag.CommandLine)

	cfg := config.GetConfig()

    initializeLogger(cfg.LogLevel)

	zap.L().Info("STDOUT Global Logger started",
		zap.String("nodeID", cfg.NodeID), 
	)

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

func initializeLogger(logLevel string) *zap.Logger {
    var zapLogLevel zapcore.Level
	logLevel = strings.ToUpper(logLevel);
	fmt.Println("Log Level:", logLevel)

    // Convertir el nivel de log a un formato válido para Zap
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
        panic("Invalid log Level!")
    }

    zapCfg := zap.Config{
        Encoding:         "json",  // Los logs serán en formato JSON
        Level:            zap.NewAtomicLevelAt(zapLogLevel),
        OutputPaths:      []string{"stdout"},  // Salida en stdout
        ErrorOutputPaths: []string{"stderr"},  // Errores a stderr
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

    logger, err := zapCfg.Build()
    if err != nil {
        panic("Failed to initialize logger: " + err.Error())
    }

    // Reemplazar el logger global
    zap.ReplaceGlobals(logger)
    return logger
}

