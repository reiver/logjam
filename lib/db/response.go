package db

import (
	"encoding/json"
	"fmt"
)

type Response struct {
	Page       int      `json:"page"`
	PerPage    int      `json:"perPage"`
	TotalItems int      `json:"totalItems"`
	TotalPages int      `json:"totalPages"`
	Items      []Record `json:"items"`
}

func parseResponse(ctx Context, jsonData string, roomName string) (*Record, error) {

	ctx.Debug("{parseResponse} BEGIN")
	defer ctx.Debug("{parseResponse} END")

	var response Response
	err := json.Unmarshal([]byte(jsonData), &response)
	if err != nil {
		ctx.Error("Error parsing JSON:", err)
		return nil, fmt.Errorf("error parsing JSON: %w", err)
	}

	for _, item := range response.Items {
		if item.Name == roomName {
			ctx.Infof("Found record: %+v", item)
			return &item, nil
		}
	}

	return nil, fmt.Errorf("record not found")
}
