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

var uselessMap map[string]any

type pocketBaseDBService struct {
	BaseURL    string
	Client     *http.Client
	AdminToken string
}

func NewPocketBaseDBService(baseURL, token string) IDBService {
	//token, err := authenticate(baseURL, username, password)
	//if err != nil {
	//	panic(err)
	//	return nil
	//}

	return &pocketBaseDBService{
		BaseURL:    baseURL,
		Client:     &http.Client{},
		AdminToken: token,
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

	return p.post(url, payload, &uselessMap)
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
		"name":   cname,
		"type":   "base",
		"fields": []any{
			/*
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
			*/
		},
		"listRule":   "",
		"viewRule":   "",
		"createRule": "",
		"updateRule": "",
		"deleteRule": "",
		"options":    map[string]any{},
	}

	getType := func(t TDBFieldType) string {
		switch t {
		case TextType:
			return "text"
		case EmailType:
			return "email"
		case UnsignedInteger64Type, FloatType, IntegerType:
			return "number"
		case BooleanType:
			return "bool"
		case DateType:
			return "date"
		case AutodateType:
			return "autodate"
		case FileType:
			return string(FileType)
		default:
			return "text"
		}
	}
	for _, f := range fields {
		payload["fields"] = append(payload["fields"].([]any), map[string]any{
			//"id":       f.Name,
			"name":     f.Name,
			"type":     getType(f.Type),
			"required": false,
			"unique":   f.Unique,
		})
	}
	payload["fields"] = append(payload["fields"].([]any),
		// created @ record creation
		map[string]any{
			"id":       "created",
			"name":     "created",
			"type":     getType(AutodateType),
			"required": false,
			"unique":   false,
			"onCreate": true,
			// PocketBase will auto‑set on create :contentReference[oaicite:0]{index=0}
		},
		// updated on every record update
		map[string]any{
			"id":       "updated",
			"name":     "updated",
			"type":     getType(AutodateType),
			"required": false,
			"unique":   false,
			"onUpdate": true,
			// PocketBase will auto‑set on update :contentReference[oaicite:1]{index=1}
		},
	)
	return p.post(url, payload, &uselessMap)
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

		if strings.Contains(k, "[") {
			field := k[:strings.Index(k, "[")]
			op := k[strings.Index(k, "[")+1 : strings.Index(k, "]")]
			symbol := map[string]string{
				"gte": ">=",
				"lte": "<=",
				"gt":  ">",
				"lt":  "<",
				"ne":  "!=",
			}[op]
			if symbol == "" {
				symbol = "="
			}
			parts = append(parts, fmt.Sprintf("(%s %s %s)", field, symbol, lit))
		} else {
			parts = append(parts, fmt.Sprintf("(%s = %s)", k, lit))
		}
		//parts = append(parts, fmt.Sprintf("(%s=%s)", k, lit))
	}
	expr := strings.Join(parts, " && ")
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

func (p *pocketBaseDBService) Insert(cname string, data map[string]any) (string, error) {
	// Build URL for creating a new record
	url := fmt.Sprintf("%s/api/collections/%s/records", p.BaseURL, cname)

	// Send POST request with the data
	var response Record
	err := p.post(url, data, &response)
	if err != nil {
		return "", err
	}

	// Return the created record's ID
	return response["id"].(string), nil
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
func (p *pocketBaseDBService) DeleteByFilter(cname string, filter map[string]any) error {
	records, err := p.GetByFilter(cname, filter)
	if err != nil {
		return err
	}

	for _, rec := range records {
		id, ok := rec["id"].(string)
		if !ok || id == "" {
			continue
		}
		if err := p.Delete(cname, id); err != nil {
			return err
		}
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

func (p *pocketBaseDBService) post(url string, body map[string]any, target any) error {
	buf, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", url, bytes.NewReader(buf))
	p.addAuth(req)
	req.Header.Set("Content-Type", "application/json")
	res, err := p.Client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	b, _ := io.ReadAll(res.Body)
	if res.StatusCode >= 300 {
		return errors.New(string(b))
	}
	return json.Unmarshal(b, target)
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
