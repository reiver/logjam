package rooms

import (
	"github.com/reiver/logjam/lib/members"
)

type MemberDTO = libmembers.DTO

type UserMessageModel struct {
	Message  string `json:"message"`
	SenderId uint64 `json:"senderId"`
}
