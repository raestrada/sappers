package config

import (
    "sync"
    "github.com/spf13/viper"
)

type Config struct {
    GossipPort int
    RaftAddr   string
    HTTPAddr   string
    NodeID     string
    Peers      []string
    LogLevel   string
	RaftDir    string
}

var (
    config *Config
    once   sync.Once
)

func GetConfig() *Config {
    once.Do(func() {
		viper.SetEnvPrefix("SAPPERS")  // Establece el prefijo para las variables de entorno
        viper.AutomaticEnv()            
        viper.SetDefault("gossip-port", 7946)
        viper.SetDefault("raft-addr", ":12000")
        viper.SetDefault("http-addr", ":11000")
        viper.SetDefault("node-id", "default-node")
        viper.SetDefault("log-level", "ERROR")  
        viper.SetDefault("peers", []string{"127.0.0.1"})
		viper.SetDefault("raft-dir", "raft/node")

        viper.BindEnv("gossip-port")
        viper.BindEnv("raft-addr")
        viper.BindEnv("http-addr")
        viper.BindEnv("node-id")
        viper.BindEnv("log-level")
        viper.BindEnv("peers")
		viper.BindEnv("raft-dir")

        // Parsear peers como una lista
        peers := viper.GetStringSlice("peers")

        config = &Config{
            GossipPort: viper.GetInt("gossip-port"),
            RaftAddr:   viper.GetString("raft-addr"),
            HTTPAddr:   viper.GetString("http-addr"),
            NodeID:     viper.GetString("node-id"),
            LogLevel:   viper.GetString("log-level"), 
            Peers:      peers,
			RaftDir:    viper.GetString("raft-dir"), 
        }
    })
    return config
}
