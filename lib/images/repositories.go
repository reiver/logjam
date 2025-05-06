package images

type IImagesRepository interface {
	Create(CreateImageDTO) error
	Get(id string) (ImageDTO, error)
	Delete(id string) error
}
