package neynar

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/reiver/logjam/lib/marshal"
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
	OwnerIdKey    = "ownerId"
)

func NewHTTPRepository(baseURL, apiKey string) INeynarServiceRepository {
	return &httpRepository{
		client:  &http.Client{Timeout: 10 * time.Second},
		apiKey:  apiKey,
		baseURL: baseURL,
	}
}

func (repo *httpRepository) SaveAccountKeys(account AK, ownerId string) error {
	rows, err := dbsrv.Repository.GetByFilter(neynarIdsTable, map[string]any{
		FIDKey: account.FID,
	})
	if err != nil {
		return err
	}
	if rows == nil || len(rows) == 0 {
		_, err := dbsrv.Repository.Insert(neynarIdsTable, map[string]any{
			FIDKey:        account.FID,
			SignerUUIDKey: account.SignerUUID,
			OwnerIdKey:    ownerId,
		})
		return err
	}
	return nil
}

func (repo *httpRepository) CreateCast(userId string, content CastPayload) error {
	rows, err := dbsrv.Repository.GetByFilter(neynarIdsTable, map[string]any{
		OwnerIdKey: userId,
	})
	if err != nil {
		return err
	}
	if rows == nil || len(rows) == 0 {
		return errors.New("submit account ids first")
	}
	id := NeynarIdDTO{}
	err = marshal.MapToObj(rows[0], &id)

	url := fmt.Sprintf("%s/v2/farcaster/cast", repo.baseURL)

	account := id.AK

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
