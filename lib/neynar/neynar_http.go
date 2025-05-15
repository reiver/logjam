package neynar

import (
	"bytes"
	"encoding/json"
	"fmt"
	cErrors "github.com/reiver/logjam/lib/errors"
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
)

func NewHTTPRepository(baseURL, apiKey string) INeynarServiceRepository {
	return &httpRepository{
		client:  &http.Client{Timeout: 10 * time.Second},
		apiKey:  apiKey,
		baseURL: baseURL,
	}
}

func (repo *httpRepository) NeynarAccountExists(fid int64) (bool, error) {
	rows, err := dbsrv.Repository.GetByFilter(neynarIdsTable, map[string]any{
		FIDKey: fid,
	})
	if err != nil {
		return false, err
	}
	if rows == nil || len(rows) == 0 {
		return false, nil
	}
	return true, nil
}

func (repo *httpRepository) SaveAccountKeys(input AK) error {
	rows, err := dbsrv.Repository.GetByFilter(neynarIdsTable, map[string]any{
		FIDKey: input.FID,
	})
	if err != nil {
		return err
	}
	if rows == nil || len(rows) == 0 {
		_, err = dbsrv.Repository.Insert(neynarIdsTable, map[string]any{
			FIDKey:        input.FID,
			SignerUUIDKey: input.SignerUUID,
		})
		if err != nil {
			return err
		}

		return nil
	} else {
		d := AK{}
		err = marshal.MapToObj(rows[0], &d)
		if err != nil {
			return err
		}
		err = dbsrv.Repository.Update(neynarIdsTable, rows[0]["id"].(string), map[string]any{
			FIDKey:        input.FID,
			SignerUUIDKey: input.SignerUUID,
		})
		if err != nil {
			return err
		}

		return nil
	}
}

func (repo *httpRepository) VerifySigner(signerUUID string) (ok bool, err error) {
	url := fmt.Sprintf("%s/v2/farcaster/signer?signer_uuid="+signerUUID, repo.baseURL)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return false, cErrors.NewErrorFromErr(http.StatusInternalServerError, fmt.Errorf("failed to create request: %w", err))
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", repo.apiKey)

	resp, err := repo.client.Do(req)
	if err != nil {
		return false, cErrors.NewErrorFromErr(http.StatusInternalServerError, fmt.Errorf("request failed: %w", err))
	}
	defer resp.Body.Close()

	if resp.StatusCode > 204 {
		respBody, _ := io.ReadAll(resp.Body)
		return false, cErrors.NewErrorFromErr(http.StatusInternalServerError, fmt.Errorf("verify failed (%d): %s", resp.StatusCode, respBody))
	}

	return true, nil
}

func (repo *httpRepository) CreateCast(fid int64, content CastPayload, signerUUID string) error {
	filter := map[string]any{
		FIDKey: fid,
	}
	if len(signerUUID) > 0 {
		filter[SignerUUIDKey] = signerUUID
	}
	rows, err := dbsrv.Repository.GetByFilter(neynarIdsTable, filter)
	if err != nil {
		return err
	}
	if rows == nil || len(rows) == 0 {
		return cErrors.NewErrorWithMsg(http.StatusUnauthorized, "submit account ids first")
	}
	id := AK{}
	err = marshal.MapToObj(rows[0], &id)

	url := fmt.Sprintf("%s/v2/farcaster/cast", repo.baseURL)

	account := id

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
		return cErrors.NewErrorFromErr(http.StatusInternalServerError, fmt.Errorf("failed to marshal cast payload: %w", err))
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return cErrors.NewErrorFromErr(http.StatusInternalServerError, fmt.Errorf("failed to create request: %w", err))
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", repo.apiKey)

	resp, err := repo.client.Do(req)
	if err != nil {
		return cErrors.NewErrorFromErr(http.StatusInternalServerError, fmt.Errorf("request failed: %w", err))
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return cErrors.NewErrorFromErr(http.StatusInternalServerError, fmt.Errorf("cast failed (%d): %s", resp.StatusCode, respBody))
	}

	return nil
}
