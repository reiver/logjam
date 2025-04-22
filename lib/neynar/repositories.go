package neynar

type AK struct {
	SignerUUID string `json:"signer_uuid"`
	FID        uint64 `json:"fid"`
}

type CastPayload struct {
	Text      string `json:"text"`
	ParentURL string `json:"parent_url,omitempty"`
	Embeds    []any  `json:"embeds,omitempty"`
}

type INeynarServiceRepository interface {
	SaveAccountKeys(account AK) error
	CreateCast(FID uint64, payload CastPayload) error
}
