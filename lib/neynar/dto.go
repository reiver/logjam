package neynar

type AK struct {
	SignerUUID string `json:"signerUUID"`
	FID        int64  `json:"fid"`
}

type CastPayload struct {
	Text      string `json:"text"`
	ParentURL string `json:"parent_url,omitempty"`
	Embeds    []any  `json:"embeds,omitempty"`
}
