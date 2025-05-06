package images

type ImageDTO struct {
	ID    string `json:"id"`
	Image []byte `json:"image"`
}

type CreateImageDTO struct {
	Image []byte `json:"image"`
}
