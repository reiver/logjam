package db

// Define structs to match the JSON structure
type Record struct {
	CollectionID   string `json:"collectionId"`
	CollectionName string `json:"collectionName"`
	Created        string `json:"created"`
	Description    string `json:"description"`
	HostID         string `json:"hostId"`
	ID             string `json:"id"`
	Name           string `json:"name"`
	Thumbnail      string `json:"thumbnail"`
	Updated        string `json:"updated"`
}
