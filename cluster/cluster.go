package cluster

import (
	"github.com/raestrada/sappers/members"
)

// Cluster ...
type Cluster struct {
	memberList members.MemberList
}

// Create ...
func Create(mfactory members.MemberListFactory) Cluster {
	return Cluster{
		memberList: mfactory.Create(),
	}
}

// Init ...
func (c Cluster) Init(peers []string) {
	c.memberList.Join(peers)
}
