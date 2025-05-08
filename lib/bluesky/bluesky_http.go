package bluesky

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	cErrors "github.com/reiver/logjam/lib/errors"
	"github.com/reiver/logjam/lib/tokens"
	"github.com/reiver/logjam/lib/users"
	dbsrv "github.com/reiver/logjam/srv/db"
	userssrv "github.com/reiver/logjam/srv/users"
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

func (repo *httpRepository) SaveLastTokens(input SubmitReqModel) (*users.CompleteSignUpResponse, error) {
	rows, err := dbsrv.Repository.GetByFilter(sessionsTable, map[string]any{
		didKey: input.DID,
	})
	if err != nil {
		return nil, err
	}
	ownerId := ""
	if rows != nil && len(rows) > 0 {
		ownerId = rows[0][ownerIdKey].(string)
		for _, r := range rows {
			dbsrv.Repository.Delete(sessionsTable, r["id"].(string))
		}
	} else {
		uid, err := userssrv.Repository.Create(users.CreateUserDTO{
			Email:    "",
			Name:     input.Name,
			UserName: "",
			Bio:      input.Bio,
		})
		if err != nil {
			return nil, err
		}
		ownerId = uid
	}

	_, err = dbsrv.Repository.Insert(sessionsTable, map[string]any{
		accessTokenKey:  input.AccessToken,
		refreshTokenKey: input.RefreshToken,
		didKey:          input.DID,
		handleKey:       input.Handle,
		ownerIdKey:      ownerId,
	})
	if err != nil {
		return nil, err
	}
	token, err := tokens.CreateToken(ownerId, nil)
	if err != nil {
		return nil, err
	}
	return &users.CompleteSignUpResponse{
		UserID: ownerId,
		Token:  token,
	}, nil
}

func (repo *httpRepository) RefreshTokens(ak AK) (AK, error) {
	url := fmt.Sprintf("%s/xrpc/com.atproto.server.refreshSession", repo.svcAddr)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return AK{}, cErrors.NewErrorFromErr(http.StatusInternalServerError, fmt.Errorf("error creating HTTP request: %w", err))
	}
	client := &http.Client{Timeout: 16 * time.Second}
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
	_, err = repo.SaveLastTokens(SubmitReqModel{AK: ak})
	return ak, err
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

func (repo *httpRepository) AccountExists(userId string) (bool, error) {
	rows, err := dbsrv.Repository.GetByFilter(sessionsTable, map[string]any{
		ownerIdKey: userId,
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
func (repo *httpRepository) CreatePost(ownerId, text string) error {
	url := fmt.Sprintf("%s/xrpc/com.atproto.repo.createRecord", repo.svcAddr)
	ak, err := repo.getAK(ownerId)
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
