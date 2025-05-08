package bluesky

import "github.com/reiver/logjam/lib/users"

type SubmitReqModel struct {
	AK
	Name string
	Bio  string
}

// AK Access Keys
type AK struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	Handle       string `json:"handle"`
	DID          string `json:"did"`
}
type TextPostRecord struct {
	Type      string `json:"$type"`
	Text      string `json:"text"`
	CreatedAt string `json:"createdAt"`
}
type IBlueSkyServiceRepository interface {
	SaveLastTokens(SubmitReqModel) (*users.CompleteSignUpResponse, error)
	RefreshTokens(accessKeys AK) (AK, error)
	CreatePost(userId, text string) error
}
