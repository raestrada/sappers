package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/spf13/viper"
	"github.com/spf13/pflag"
	"github.com/raestrada/sappers/config"
	"github.com/raestrada/sappers/cluster"
	"github.com/raestrada/sappers/members"
	"github.com/raestrada/sappers/consensus"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Command line defaults
const (
	DefaultHTTPAddr = ":11000"
	DefaultRaftAddr = ":12000"
)

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

	// Obtener la configuración desde el singleton
	cfg := config.GetConfig()

	// Inicializar el logger
	initializeLogger(cfg.LogLevel)

	zap.L().Info("STDOUT Global Logger started", zap.String("nodeID", cfg.NodeID))

	// Crear un contexto para manejar la interrupción
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Capturar señales para hacer una finalización ordenada
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Iniciar el clúster en una goroutine
	go startCluster(ctx)

	// Esperar la señal de interrupción
	select {
	case sig := <-sigs:
		fmt.Println("Signal received:", sig)
		cancel() // Cancelar el contexto para finalizar el cluster
	}

	fmt.Println("Shutting down gracefully...")
}

// startCluster inicia el clúster
func startCluster(ctx context.Context) {

	// Crear una instancia de GossipMemberListFactory
	gossipFactory := members.GossipMemberListFactory{}

	consensusFactory := consensus.ConsensusFactory{}

	// Crear el clúster
	clusterInstance := cluster.Create(gossipFactory, consensusFactory)

	// Inicializar el clúster y pasar el contexto y los peers
	clusterInstance.Init(ctx)
}

// initializeLogger inicializa el logger global
func initializeLogger(logLevel string) *zap.Logger {
	var zapLogLevel zapcore.Level
	logLevel = strings.ToUpper(logLevel)
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
		panic("Invalid log level!")
	}

	zapCfg := zap.Config{
		Encoding:         "json",  // Los logs serán en formato JSON
		Level:            zap.NewAtomicLevelAt(zapLogLevel),
		OutputPaths:      []string{"stdout"},  // Salida en stdout
		ErrorOutputPaths: []string{"stderr"},  // Errores a stderr
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "msg",
			LevelKey:   "level",
			EncodeLevel: zapcore.CapitalLevelEncoder,
			TimeKey:    "time",
			EncodeTime: zapcore.ISO8601TimeEncoder,
			CallerKey:  "caller",
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
