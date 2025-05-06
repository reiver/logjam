package layouts

type CreateLayoutDTO struct {
	Name     string `json:"name"`
	OwnerId  string `json:"ownerId"`
	Data     string `json:"data"`
	Type     string `json:"type"`
	FileHash string `json:"fileHash"`
	LastUsed bool   `json:"lastUsed"`
}

type LayoutDTO struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	OwnerId  string `json:"ownerId"`
	Data     string `json:"data"`
	Type     string `json:"type"`
	FileHash string `json:"fileHash"`
	LastUsed bool   `json:"lastUsed"`
}
