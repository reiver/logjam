package bluesky

// AK Access Keys
type AK struct {
	AccessToken  string `json:"accessJwt"`
	RefreshToken string `json:"refreshJwt"`
	Handle       string `json:"handle"`
	DID          string `json:"did"`
}
type TextPostRecord struct {
	Type      string `json:"$type"`
	Text      string `json:"text"`
	CreatedAt string `json:"createdAt"`
}

type IBlueSkyServiceRepository interface {
	RefreshTokens(did string) (AK, error)
	CreatePost(did, text string) error
	AccountExists(did string) (bool, error)
}
