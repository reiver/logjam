package bluesky

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	cErrors "github.com/reiver/logjam/lib/errors"
	"github.com/reiver/logjam/lib/marshal"
	dbsrv "github.com/reiver/logjam/srv/db"
	"io"
	"net/http"
	"sync"
	"time"
)

type httpRepository struct {
	client  *http.Client
	svcAddr string
	*sync.Mutex
}

const (
	sessionsTable = "blueSkySessions"

	//keys
	accessTokenKey  = "accessJwt"
	refreshTokenKey = "refreshJwt"
	handleKey       = "handle"
	didKey          = "did"
)

func NewHTTPRepository(svcAddr string) IBlueSkyServiceRepository {
	return &httpRepository{
		client: &http.Client{
			Timeout: 8 * time.Second,
		},
		svcAddr: svcAddr,
		Mutex:   &sync.Mutex{},
	}
}

func (repo *httpRepository) RefreshTokens(did string, accessKey, refreshKey string) (AK, error) {
	url := fmt.Sprintf("%s/xrpc/com.atproto.server.refreshSession", repo.svcAddr)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return AK{}, cErrors.NewErrorFromErr(http.StatusInternalServerError, fmt.Errorf("error creating HTTP request: %w", err))
	}
	client := &http.Client{Timeout: 16 * time.Second}

	filters := map[string]any{didKey: did}
	if len(accessKey) > 0 && len(refreshKey) > 0 {
		filters[accessTokenKey] = accessKey
		filters[refreshTokenKey] = refreshKey
	}
	rows, err := dbsrv.Repository.GetByFilter(sessionsTable, filters)
	if err != nil {
		return AK{}, err
	}
	if len(rows) == 0 {
		return AK{}, cErrors.NewErrorWithMsg(http.StatusNotFound, "couldnt find tokens for this account did/at/rt")
	}
	var ak AK
	err = marshal.MapToObj(rows[0], &ak)
	if err != nil {
		return AK{}, nil
	}
	req.Header.Set("Authorization", "Bearer "+ak.RefreshToken)
	resp, err := client.Do(req)
	if err != nil {
		return AK{}, cErrors.NewErrorFromErr(http.StatusInternalServerError, fmt.Errorf("error sending HTTP request: %w", err))
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return AK{}, cErrors.NewErrorFromErr(http.StatusInternalServerError, fmt.Errorf("refresh token request failed with status %d: %s", resp.StatusCode, string(body)))
	}

	// Parse the response.
	var res struct {
		DID        string         `json:"did"`
		DidDoc     map[string]any `json:"didDoc"`
		Handle     string         `json:"handle"`
		AccessJWT  string         `json:"accessJwt"`
		RefreshJWT string         `json:"refreshJwt"`
		Active     bool           `json:"active"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return AK{}, cErrors.NewErrorFromErr(http.StatusInternalServerError, fmt.Errorf("error decoding refresh response: %w", err))
	}

	// Verify the response contains new tokens.
	if res.AccessJWT == "" || res.RefreshJWT == "" {
		return AK{}, cErrors.NewErrorFromErr(http.StatusInternalServerError, errors.New("invalid refresh response: missing tokens"))
	}

	ak = AK{
		AccessToken:  res.AccessJWT,
		RefreshToken: res.RefreshJWT,
		Handle:       res.Handle,
		DID:          res.DID,
	}
	err = dbsrv.Repository.UpdateByFilter(sessionsTable, map[string]any{didKey: ak.DID}, map[string]any{
		accessTokenKey:  res.AccessJWT,
		refreshTokenKey: res.RefreshJWT,
	})
	return ak, err
}

func (repo *httpRepository) getAK(did string) (*AK, error) {
	recs, err := dbsrv.Repository.GetByFilter(sessionsTable, map[string]any{didKey: did})
	if err != nil {
		return nil, err
	}
	if len(recs) == 0 {
		return nil, nil
	}
	r0 := recs[0]
	return &AK{
		AccessToken:  r0[accessTokenKey].(string),
		RefreshToken: r0[refreshTokenKey].(string),
		Handle:       r0[handleKey].(string),
		DID:          r0[didKey].(string),
	}, nil
}

func (repo *httpRepository) AccountExists(did string) (bool, error) {
	rows, err := dbsrv.Repository.GetByFilter(sessionsTable, map[string]any{
		didKey: did,
	})
	if err != nil {
		return false, err
	}
	if rows != nil && len(rows) == 0 {
		return false, nil
	}
	return true, nil
}

// CreatePost creates a new post by calling the create post endpoint.
func (repo *httpRepository) CreatePost(did, text string) error {
	url := fmt.Sprintf("%s/xrpc/com.atproto.repo.createRecord", repo.svcAddr)
	ak, err := repo.getAK(did)
	if err != nil {
		return err
	}
	if ak == nil {
		return cErrors.NewErrorFromErr(http.StatusUnauthorized, errors.New("couldn't find account keys, maybe not submitted yet"))
	}
	post := TextPostRecord{
		Type:      "app.bsky.feed.post",
		Text:      text,
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
	payload := struct {
		Repo       string      `json:"repo"`
		Collection string      `json:"collection"`
		Record     interface{} `json:"record"`
	}{
		Repo:       ak.DID,
		Collection: "app.bsky.feed.post",
		Record:     post,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return cErrors.NewErrorFromErr(http.StatusInternalServerError, fmt.Errorf("error marshaling payload: %w", err))
	}

	// Prepare the request.
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return cErrors.NewErrorFromErr(http.StatusInternalServerError, fmt.Errorf("error creating HTTP request: %w", err))
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+ak.AccessToken)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return cErrors.NewErrorFromErr(http.StatusInternalServerError, fmt.Errorf("error sending HTTP request: %w", err))
	}
	defer resp.Body.Close()

	// Check for non-200 status codes.
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return cErrors.NewErrorFromErr(resp.StatusCode, fmt.Errorf("create post failed (status %d): %s", resp.StatusCode, string(body)))
	}

	// Optionally, decode and use the response data here.
	return nil
}
