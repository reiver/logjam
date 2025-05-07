package bluesky

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
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
	accessTokenKey  = "accessToken"
	refreshTokenKey = "refreshToken"
	handleKey       = "handle"
	didKey          = "did"
	ownerIdKey      = "ownerId"
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

func (repo *httpRepository) SaveLastTokens(accessKeys AK, ownerId string) error {
	rows, err := dbsrv.Repository.GetByFilter(sessionsTable, map[string]any{
		ownerIdKey: ownerId,
	})
	if err != nil {
		return err
	}
	if rows != nil && len(rows) > 0 {
		for _, r := range rows {
			dbsrv.Repository.Delete(sessionsTable, r["id"].(string))
		}
	}
	_, err = dbsrv.Repository.Insert(sessionsTable, map[string]any{
		accessTokenKey:  accessKeys.AccessToken,
		refreshTokenKey: accessKeys.RefreshToken,
		didKey:          accessKeys.DID,
		handleKey:       accessKeys.Handle,
		ownerIdKey:      ownerId,
	})
	return err
}

func (repo *httpRepository) RefreshTokens(ak AK) (AK, error) {
	url := fmt.Sprintf("%s/xrpc/com.atproto.server.refreshSession", repo.svcAddr)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return AK{}, fmt.Errorf("error creating HTTP request: %w", err)
	}
	client := &http.Client{Timeout: 16 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return AK{}, fmt.Errorf("error sending HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return AK{}, fmt.Errorf("refresh token request failed with status %d: %s", resp.StatusCode, string(body))
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
		return AK{}, fmt.Errorf("error decoding refresh response: %w", err)
	}

	// Verify the response contains new tokens.
	if res.AccessJWT == "" || res.RefreshJWT == "" {
		return AK{}, errors.New("invalid refresh response: missing tokens")
	}

	ak = AK{
		AccessToken:  res.AccessJWT,
		RefreshToken: res.RefreshJWT,
		Handle:       res.Handle,
		DID:          res.DID,
	}
	return ak, repo.SaveLastTokens(ak, "")
}

func (repo *httpRepository) getAK(ownerId string) (*AK, error) {
	recs, err := dbsrv.Repository.GetByFilter(sessionsTable, map[string]any{ownerIdKey: ownerId})
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

// CreatePost creates a new post by calling the create post endpoint.
func (repo *httpRepository) CreatePost(ownerId, text string) error {
	url := fmt.Sprintf("%s/xrpc/com.atproto.repo.createRecord", repo.svcAddr)
	ak, err := repo.getAK(ownerId)
	if err != nil {
		return err
	}
	if ak == nil {
		return errors.New("couldn't find account keys, maybe not submitted yet")
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
		return fmt.Errorf("error marshaling payload: %w", err)
	}

	// Prepare the request.
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Errorf("error creating HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+ak.AccessToken)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending HTTP request: %w", err)
	}
	defer resp.Body.Close()

	// Check for non-200 status codes.
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("create post failed (status %d): %s", resp.StatusCode, string(body))
	}

	// Optionally, decode and use the response data here.
	return nil
}
