package members

import (
	"github.com/raestrada/sappers/config"
	"github.com/hashicorp/memberlist"
	"go.uber.org/zap"
)

// MemberListFactory defines the factory interface for creating member lists.
type MemberListFactory interface {
	Create() MemberList
}

// Join hace que este nodo se una a un cluster utilizando los peers proporcionados.
func (mla *MemberlistAdapter) Join(peers []string) error {
	_, err := mla.list.Join(peers)
	if err != nil {
		zap.L().Fatal(
			"Failed to join cluster",
			zap.String("type", "Join"),
			zap.String("msg", err.Error()),
		)
		return err
	}
	zap.L().Info("Successfully joined the cluster")
	return nil
}

// Get retorna la lista de miembros conectados.
func (mla *MemberlistAdapter) Get() []Member {
	members := make([]Member, len(mla.list.Members()))
	for i, member := range mla.list.Members() {
		members[i] = Member{
			Addr: member.Addr.String(),
			Name: member.Name,
		}
	}
	return members
}

// Create crea una nueva instancia de MemberlistAdapter utilizando memberlist.
func (mf MemberlistFactory) Create() MemberList {
	cfg := config.GetConfig()
	mlConfig := memberlist.DefaultLocalConfig()
	mlConfig.BindPort = cfg.GossipPort
	mlConfig.Name = cfg.NodeID

	list, err := memberlist.Create(mlConfig)
	if err != nil {
		zap.L().Fatal(
			"Failed to create memberlist",
			zap.String("type", "Create"),
			zap.String("msg", err.Error()),
		)
	}

	zap.L().Info("Memberlist created successfully",
		zap.String("nodeID", cfg.NodeID),
		zap.Int("gossipPort", cfg.GossipPort),
	)

	return &MemberlistAdapter{
		list: list,
	}
}
