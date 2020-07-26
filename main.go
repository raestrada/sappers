package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"

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

func main() {
	flag.BoolVar(&inmem, "inmem", false, "Use in-memory storage for Raft")
	flag.StringVar(&httpAddr, "haddr", DefaultHTTPAddr, "Set the HTTP bind address")
	flag.StringVar(&raftAddr, "raddr", DefaultRaftAddr, "Set Raft bind address")
	flag.StringVar(&joinAddr, "join", "", "Set join address, if any")
	flag.StringVar(&nodeID, "id", "", "Node ID")
	flag.StringVar(&inPeers, "peers", "127.0.0.1", "peers list")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] <raft-data-path> \n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	if flag.NArg() == 0 {
		fmt.Fprintf(os.Stderr, "No Raft storage directory specified\n")
		os.Exit(1)
	}

	peers = strings.Split(inPeers, ",")

	if value, ok := os.LookupEnv("SAPPERS_PEERS"); ok {
		peers = strings.Split(value, ",")
	}

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
	var cluster = cluster.Create(members.GossipMemberListFactory{})
	cluster.Init(peers)
}
