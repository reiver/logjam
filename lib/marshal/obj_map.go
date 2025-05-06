package marshal

import (
	"encoding/json"
	"github.com/reiver/logjam/lib/db"
)

func ObjToMap(obj any) (result map[string]any, err error) {
	b, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, &result)
	return
}
func MapToObj(data map[string]any, out any) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, out)
}

func MapArrayToObjArray[I map[string]any | db.Record, O any](data []I, out *[]O) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, out)
}
