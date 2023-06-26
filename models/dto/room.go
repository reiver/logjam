package dto

type MemberDTO struct {
	Id       uint64 `json:"id"`
	Name     string `json:"name"`
	Role     string `json:"role"`
	StreamId string `json:"streamId"`
}
