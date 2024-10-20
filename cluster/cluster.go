package cluster

import (
	"context"
	"fmt"

	"github.com/raestrada/sappers/config"
	"github.com/raestrada/sappers/consensus"
	"github.com/raestrada/sappers/members"
	"go.uber.org/zap"
)

// Cluster manages the lifecycle of a distributed cluster.
type Cluster struct {
	memberList members.MemberList
	consensus  *consensus.Consensus
}

// Create initializes a new cluster with the provided member list factory and consensus instance.
func Create(mfactory members.MemberListFactory, consensusFactory consensus.ConsensusFactory) Cluster {
	fmt.Println("---------------------------------------------")
	fmt.Println(" _____                                        ")
	fmt.Println("/  ___|                                      ")
	fmt.Println("\\ `--.  __ _ _ __  _ __   ___ _ __ ___       ")
	fmt.Println(" `--. \\/ _` | '_ \\| '_ \\ / _ \\ '__/ __|      ")
	fmt.Println(" /\\__/ / (_| | |_) | |_) |  __/ |  \\__ \\     ")
	fmt.Println(" \\____/ \\__,_| .__/| .__/ \\___|_|  |___/     ")
	fmt.Println("             | |   | |                       ")
	fmt.Println("             |_|   |_|                       ")
	fmt.Println("---------------------------------------------")
	fmt.Println("     Who is in charge here, Commander?       ")
	fmt.Println("---------------------------------------------")

	zap.L().Info("Creating Cluster ...")

	// Return a Cluster instance with the member list and consensus logic
	mlist := mfactory.Create()
	return Cluster{
		memberList: mlist,
		consensus: consensusFactory.Create(mlist),
	}
}

// Init initializes the cluster, joins peers, and starts the consensus process.
func (c Cluster) Init(ctx context.Context) {
	zap.L().Info("Starting Cluster ...")
	cfg := config.GetConfig()

	// Join the cluster using gossip
	c.memberList.Join(cfg.Peers)

	// Initialize the consensus mechanism (Raft)
	zap.L().Info("Initializing consensus mechanism ...")
	c.consensus.Init(ctx)

	zap.L().Info("Cluster successfully started and consensus mechanism initialized.")
}
