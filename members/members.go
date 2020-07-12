package members

import (
  "github.com/raestrada/sappers/domain"
)

type MemberListFactory struct {}

type MemberList interface {
  Join(peers)
  Get() []domain.Member
}
