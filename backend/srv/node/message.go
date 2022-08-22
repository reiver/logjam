package node

import (
	"encoding/json"
	"fmt"
)

type MessageContract[T MessageData] struct {
	Type string `json:"type"`
	Data T      `json:"data"`
}

type MessageData interface {
	JoinMessageData | ConnectMessageData
}

type JoinMessageData struct {
	Name string `json:"name"`
	Role string `json:"role"`
}

type ConnectMessageData struct {
	Target string `json:"target"`
}

func (m *MessageContract[T]) handle() (*MessageContract[T], error) {
	fmt.Println("Handle ", m.Type)
	switch m.Type {
	case "join":

	}

	return m, nil
}

func UnmarshalJSON(data []byte) (interface{}, error) {
	var tmp *map[string]interface{} = &map[string]interface{}{}
	err := json.Unmarshal(data, tmp)
	if err != nil {
		return nil, err
	}
	if (*tmp)["type"] == "join" {
		var data JoinMessageData = JoinMessageData{
			Name: ((*tmp)["data"]).(map[string]interface{})["name"].(string),
			Role: ((*tmp)["data"]).(map[string]interface{})["role"].(string),
		}
		return NewMessageContract((*tmp)["type"].(string), data), nil
	}
	return nil, nil
}

func NewMessageContract[T MessageData](type_ string, data T) *MessageContract[T] {
	return &MessageContract[T]{
		Type: type_,
		Data: data,
	}
}
