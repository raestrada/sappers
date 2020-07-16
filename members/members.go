package members

import (
  "github.com/raestrada/sappers/domain"
)

// MemberListFactory ...
type MemberListFactory struct {}

// MemberList ...
type MemberList interface {
  Join(peers []string)
  Get() []domain.Member
}

// MemberListFactory ...
type MemberListFactory interface {
	Create() MemberList
}
