package neynar

type NeynarIdDTO struct {
	AK
	OwnerId string `json:"ownerId"`
}

type SubmitReqModel struct {
	AK
	Name string `json:"name"`
	Bio  string `json:"bio"`
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
