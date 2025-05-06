package layouts

type ILayoutsRepository interface {
	Create(input CreateLayoutDTO) (string, error)
	GetUserLayouts(userId string) ([]LayoutDTO, error)
	DeleteLayout(id, ownerId string) error
}
