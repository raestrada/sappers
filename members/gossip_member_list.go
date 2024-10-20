package members

import (
	"github.com/raestrada/sappers/config"
	"github.com/hashicorp/memberlist"
	"github.com/raestrada/sappers/domain"
	"github.com/wesovilabs/koazee"
	"go.uber.org/zap"
)

// GossipMemberList ...
type GossipMemberList struct {
	list *memberlist.Memberlist
}

// Join ..
func (gml *GossipMemberList) Join(peers []string) {
	var funcDesc = "GossipMemberList - Join"

	_, err := gml.list.Join(peers)
	if err != nil {
		zap.L().Fatal(
			funcDesc,
			zap.String("type", "Failed to join cluster"),
			zap.String("msg", err.Error()),
		)
	}
}

// Get ...
func (gml *GossipMemberList) Get() []domain.Member {
	return koazee.StreamOf(gml.list.Members()).
		Map(
			func(member domain.Member) domain.Member {
				return domain.Member{
					Addr: member.Addr,
					Name: member.Name,
				}
			}).Out().Val().([]domain.Member)
}

// GossipMemberListFactory ...
type GossipMemberListFactory struct{}

// GossipMemberListFactory - Create
func (gmlf GossipMemberListFactory) Create() MemberList {
    var funcDesc = "GossipMemberListFactory - Create"

    cfg := config.GetConfig()  // Obtener la configuración desde el singleton
    config := memberlist.DefaultLocalConfig()
    config.BindPort = cfg.GossipPort  // Usamos el puerto gossip definido en la configuración

    list, err := memberlist.Create(config)
    if err != nil {
        zap.L().Fatal(  // Usar el logger global ya configurado
            funcDesc,
            zap.String("type", "Failed to create memberlist"),
            zap.String("msg", err.Error()),
        )
    }

    zap.L().Info("Memberlist created successfully",  // Usar el logger global
        zap.String("nodeID", cfg.NodeID),
        zap.Int("gossipPort", cfg.GossipPort),
    )

    return &GossipMemberList{
        list: list,
    }
}
