package db

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Record map[string]any

type pocketBaseDBService struct {
	BaseURL    string
	Client     *http.Client
	AdminToken string
}

func NewPocketBaseDBService(baseURL, adminToken string) IDBService {
	return &pocketBaseDBService{
		BaseURL:    baseURL,
		Client:     &http.Client{},
		AdminToken: adminToken,
	}
}

func (p *pocketBaseDBService) Ping() error {
	url := fmt.Sprintf("%s/api/health", p.BaseURL)
	req, _ := http.NewRequest("GET", url, nil)
	p.addAuth(req)

	resp, err := p.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("ping failed: %s", string(body))
	}
	return nil
}
func (p *pocketBaseDBService) UpdateByFilter(cname string, filter, data map[string]any) error {
	records, err := p.GetByFilter(cname, filter)
	if err != nil {
		return err
	}

	for _, rec := range records {
		idVal, ok := rec["id"].(string)
		if !ok || idVal == "" {
			continue
		}
		if err := p.Update(cname, idVal, data); err != nil {
			return err
		}
	}

	return nil
}

func (p *pocketBaseDBService) Initialize() error {
	// skip if already there
	cols, err := p.getCollections()
	if err != nil {
		return err
	}
	for _, c := range cols {
		if c["name"] == "people" {
			return nil
		}
	}

	// define the people schema
	url := fmt.Sprintf("%s/api/collections", p.BaseURL)
	payload := map[string]any{
		"name": "people",
		"type": "base",
		"schema": []any{
			map[string]any{
				"id":       "name",
				"name":     "name",
				"type":     "text",
				"required": true,
				"unique":   false,
				"options": map[string]any{
					"min":     1,
					"max":     0,
					"pattern": "",
				},
			},
			map[string]any{
				"id":       "age",
				"name":     "age",
				"type":     "number",
				"required": true,
				"unique":   false,
				"options": map[string]any{
					"min": 0,
					"max": 0,
				},
			},
		},
		"listRule":   "",
		"viewRule":   "",
		"createRule": "",
		"updateRule": "",
		"deleteRule": "",
		"options":    map[string]any{},
	}

	return p.post(url, payload)
}
func (p *pocketBaseDBService) CreateTableIfNotExists(cname string, fields []Field) error {
	// check existing
	cols, err := p.getCollections()
	if err != nil {
		return err
	}
	for _, c := range cols {
		if c["name"] == cname {
			return nil
		}
	}

	// create with minimal required fields
	url := fmt.Sprintf("%s/api/collections", p.BaseURL)
	payload := map[string]any{
		"name":       cname,
		"type":       "base",
		"schema":     []any{}, // add fields later via Update
		"listRule":   "",
		"viewRule":   "",
		"createRule": "",
		"updateRule": "",
		"deleteRule": "",
		"options":    map[string]any{},
	}
	return p.post(url, payload)
}

func (p *pocketBaseDBService) GetById(cname, id string) (Record, error) {
	url := fmt.Sprintf("%s/api/collections/%s/records/%s", p.BaseURL, cname, id)
	var rec Record
	if err := p.get(url, &rec); err != nil {
		return nil, err
	}
	return rec, nil
}

// GetByFilter fetches records matching simple equality filters.
func (p *pocketBaseDBService) GetByFilter(cname string, filters map[string]any) ([]Record, error) {
	// build PocketBase filter expr: (field="val")&&(age=30)
	parts := make([]string, 0, len(filters))
	for k, v := range filters {
		var lit string
		switch t := v.(type) {
		case string:
			lit = fmt.Sprintf(`"%s"`, t)
		default:
			lit = fmt.Sprintf("%v", t)
		}
		parts = append(parts, fmt.Sprintf("(%s=%s)", k, lit))
	}
	expr := strings.Join(parts, "&&")
	u := fmt.Sprintf("%s/api/collections/%s/records?filter=%s", p.BaseURL, cname, url.QueryEscape(expr))

	var resp struct {
		Items []Record `json:"items"`
	}
	if err := p.get(u, &resp); err != nil {
		return nil, err
	}
	return resp.Items, nil
}

func (p *pocketBaseDBService) GetAll(cname string) ([]Record, error) {
	url := fmt.Sprintf("%s/api/collections/%s/records", p.BaseURL, cname)
	var resp struct {
		Items []Record `json:"items"`
	}
	if err := p.get(url, &resp); err != nil {
		return nil, err
	}
	return resp.Items, nil
}

func (p *pocketBaseDBService) Update(cname, id string, data map[string]any) error {
	url := fmt.Sprintf("%s/api/collections/%s/records/%s", p.BaseURL, cname, id)
	return p.patch(url, data)
}

func (p *pocketBaseDBService) Delete(cname, id string) error {
	url := fmt.Sprintf("%s/api/collections/%s/records/%s", p.BaseURL, cname, id)
	req, _ := http.NewRequest("DELETE", url, nil)
	p.addAuth(req)
	res, err := p.Client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode >= 300 {
		b, _ := io.ReadAll(res.Body)
		return errors.New(string(b))
	}
	return nil
}

// —— helpers ——

func (p *pocketBaseDBService) getCollections() ([]Record, error) {
	url := fmt.Sprintf("%s/api/collections", p.BaseURL)
	var resp struct {
		Items []Record `json:"items"`
	}
	if err := p.get(url, &resp); err != nil {
		return nil, err
	}
	return resp.Items, nil
}

func (p *pocketBaseDBService) get(url string, target any) error {
	req, _ := http.NewRequest("GET", url, nil)
	p.addAuth(req)
	res, err := p.Client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return parseResp(res, target)
}

func (p *pocketBaseDBService) post(url string, body map[string]any) error {
	buf, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", url, bytes.NewReader(buf))
	p.addAuth(req)
	req.Header.Set("Content-Type", "application/json")
	res, err := p.Client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return checkResp(res)
}

func (p *pocketBaseDBService) patch(url string, body map[string]any) error {
	buf, _ := json.Marshal(body)
	req, _ := http.NewRequest("PATCH", url, bytes.NewReader(buf))
	p.addAuth(req)
	req.Header.Set("Content-Type", "application/json")
	res, err := p.Client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return checkResp(res)
}

func (p *pocketBaseDBService) addAuth(req *http.Request) {
	if p.AdminToken != "" {
		req.Header.Set("Authorization", "Bearer "+p.AdminToken)
	}
}

func parseResp(res *http.Response, target any) error {
	b, _ := io.ReadAll(res.Body)
	if res.StatusCode >= 300 {
		return errors.New(string(b))
	}
	return json.Unmarshal(b, target)
}

func checkResp(res *http.Response) error {
	b, _ := io.ReadAll(res.Body)
	if res.StatusCode >= 300 {
		return errors.New(string(b))
	}
	return nil
}
