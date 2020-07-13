package members

import (
  "github.com/raestrada/sappers/domain"
)

// MemberListFactory ...
type MemberListFactory struct {}

// MemberList ...
type MemberList interface {
  Join(peers)
  Get() []domain.Member
}
