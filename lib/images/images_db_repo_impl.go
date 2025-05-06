package images

type imageRepo struct {
}

func GetNewImageRepository() IImagesRepository {
	return &imageRepo{}
}

func (i *imageRepo) Create(dto CreateImageDTO) error {
	//TODO implement me
	panic("implement me")
}

func (i *imageRepo) Get(id string) (ImageDTO, error) {
	//TODO implement me
	panic("implement me")
}

func (i *imageRepo) Delete(id string) error {
	//TODO implement me
	panic("implement me")
}
