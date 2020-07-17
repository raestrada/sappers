package members

import (
	"github.com/raestrada/sappers/domain"
)

// MemberList ...
type MemberList interface {
	Join(peers []string)
	Get() []domain.Member
}

// MemberListFactory ...
type MemberListFactory interface {
	Create() MemberList
}
