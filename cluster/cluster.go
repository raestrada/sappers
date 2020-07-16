package cluster

import (
  "github.com/raestrada/sappers/member"
)

// Cluster ...
type Cluster struct {
  memberList member.MemberList
}

// ClusterFactory ...
type ClusterFactory {}

// Create ...
func (cf ClusterFactory) Create(mfactory member.MemberListFactory) Cluster {
  return Cluster{
    memberList: mfactory.Create()
  }
}

// Init ...
func (c Cluster) Init(peers []string) {
  c.memberList.Join(peers)
}
