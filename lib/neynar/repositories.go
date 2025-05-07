package neynar

type NeynarIdDTO struct {
	AK
	OwnerId string `json:"ownerId"`
}
type AK struct {
	SignerUUID string `json:"signerUUID"`
	FID        uint64 `json:"fid"`
}

type CastPayload struct {
	Text      string `json:"text"`
	ParentURL string `json:"parent_url,omitempty"`
	Embeds    []any  `json:"embeds,omitempty"`
}

type INeynarServiceRepository interface {
	SaveAccountKeys(account AK, ownerId string) error
	CreateCast(userId string, payload CastPayload) error
}
