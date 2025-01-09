package rooms

type MemberDTO struct {
	Id       uint64 `json:"id"`
	Name     string `json:"name"`
	Role     string `json:"role"`
	StreamId string `json:"streamId"`
}

type UserMessageModel struct {
	Message  string `json:"message"`
	SenderId uint64 `json:"senderId"`
}
