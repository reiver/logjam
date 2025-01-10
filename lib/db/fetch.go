package db

import (
	"io/ioutil"
	"net/http"
)

func fetchRecordFromPocketBase(ctx Context, roomName string) (*Record, error) {

	ctx.Debug("{fetchRecordFromPocketBase} BEGIN")
	defer ctx.Debug("{fetchRecordFromPocketBase} END")

	url := ctx.URL("collections/rooms/records?sort=-created")

	// Create a new request using http
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		ctx.Error("Error creating request:", err)
		return nil, err
	}

	ctx.Info("Request:", req)

	// Send the request via a client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		ctx.Error("Error sending request:", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ctx.Error("Error reading response body:", err)
		return nil, err
	}

	// ctx.Info(logtag, "Response body:", body)
	record, err := parseResponse(ctx, string(body), roomName)

	return record, err
	// ctx.Info(logtag, "Response:", string(body))
}
