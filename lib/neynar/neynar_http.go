package neynar

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	dbsrv "github.com/reiver/logjam/srv/db"
	"io"
	"net/http"
	"time"
)

type httpRepository struct {
	client  *http.Client
	apiKey  string
	baseURL string
}

const (
	neynarIdsTable = "neynarIDs"

	///keys
	FIDKey        = "fid"
	SignerUUIDKey = "signerUUID"
)

func NewHTTPRepository(baseURL, apiKey string) INeynarServiceRepository {
	return &httpRepository{
		client:  &http.Client{Timeout: 10 * time.Second},
		apiKey:  apiKey,
		baseURL: baseURL,
	}
}

func (repo *httpRepository) SaveAccountKeys(account AK) error {
	return dbsrv.Repository.UpdateByFilter(neynarIdsTable, map[string]any{FIDKey: account.FID}, map[string]any{FIDKey: account.FID, SignerUUIDKey: account.SignerUUID})
}

func (repo *httpRepository) CreateCast(FID uint64, content CastPayload) error {
	url := fmt.Sprintf("%s/v2/farcaster/cast", repo.baseURL)

	account := AK{}
	recs, err := dbsrv.Repository.GetByFilter(neynarIdsTable, map[string]any{FIDKey: FID})
	if err != nil {
		return err
	}
	if len(recs) == 0 {
		return errors.New("no matching fid found, maybe ids are not submitted yet")
	}
	account.FID = recs[0][FIDKey].(uint64)
	account.SignerUUID = recs[0][SignerUUIDKey].(string)

	payload := map[string]interface{}{
		"text":        content.Text,
		"signer_uuid": account.SignerUUID,
	}

	if content.ParentURL != "" {
		payload["parent_url"] = content.ParentURL
	}
	if content.Embeds != nil {
		payload["embeds"] = content.Embeds
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal cast payload: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", repo.apiKey)

	resp, err := repo.client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("cast failed (%d): %s", resp.StatusCode, respBody)
	}

	return nil
}
