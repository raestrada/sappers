package cluster

import (
	"fmt"

	"github.com/raestrada/sappers/members"
	"go.uber.org/zap"
)

// Cluster ...
type Cluster struct {
	memberList members.MemberList
}

// Create ...
func Create(mfactory members.MemberListFactory) Cluster {
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

	return Cluster{
		memberList: mfactory.Create(),
	}
}

// Init ...
func (c Cluster) Init(peers []string) {
	zap.L().Info("Starting Cluster ...")
	c.memberList.Join(peers)
}
