package neynar

type INeynarServiceRepository interface {
	SaveAccountKeys(ak AK) error
	CreateCast(fid int64, payload CastPayload, signerUUID string) error
	NeynarAccountExists(fid int64) (bool, error)
}
