package rooms

type UserMessageModel struct {
	Message  string `json:"message"`
	SenderId uint64 `json:"senderId"`
}
