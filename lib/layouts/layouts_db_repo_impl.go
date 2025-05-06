package layouts

import (
	"github.com/reiver/logjam/lib/marshal"
	dbsrv "github.com/reiver/logjam/srv/db"
)

type layoutsRepo struct {
}

const (
	layoutsTbl = "layouts"

	ownerIdKey = "ownerId"
)

func GetNewLayoutsRepository() ILayoutsRepository {
	return &layoutsRepo{}
}

func (l *layoutsRepo) Create(input CreateLayoutDTO) (string, error) {
	data, err := marshal.ObjToMap(input)
	if err != nil {
		return "", err
	}
	return dbsrv.Repository.Insert(layoutsTbl, data)
}

func (l *layoutsRepo) GetUserLayouts(userId string) (layouts []LayoutDTO, err error) {
	result, err := dbsrv.Repository.GetByFilter(layoutsTbl, map[string]any{ownerIdKey: userId})
	if err != nil {
		return nil, err
	}
	if result == nil || len(result) == 0 {
		return nil, nil
	}
	err = marshal.MapArrayToObjArray(result, &layouts)
	return
}

func (l *layoutsRepo) DeleteLayout(id, ownerId string) error {
	//TODO implement me
	panic("implement me")
}
