package members

import (
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

// Create ...
func (gmlf GossipMemberListFactory) Create() MemberList {
	var funcDesc = "GossipMemberListFactory - Create"

	list, err := memberlist.Create(memberlist.DefaultLocalConfig())
	if err != nil {
		zap.L().Fatal(
			funcDesc,
			zap.String("type", "Failed to create memberlist"),
			zap.String("msg", err.Error()),
		)
	}
	return &GossipMemberList{
		list: list,
	}
}
