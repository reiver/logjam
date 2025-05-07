package bluesky

// AK Access Keys
type AK struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Handle       string `json:"handle"`
	DID          string `json:"did"`
}
type TextPostRecord struct {
	Type      string `json:"$type"`
	Text      string `json:"text"`
	CreatedAt string `json:"createdAt"`
}
type IBlueSkyServiceRepository interface {
	SaveLastTokens(accessKeys AK, ownerId string) error
	RefreshTokens(accessKeys AK) (AK, error)
	CreatePost(did, text string) error
}
