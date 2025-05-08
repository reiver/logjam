package neynar

import "github.com/reiver/logjam/lib/users"

type INeynarServiceRepository interface {
	SaveAccountKeys(SubmitReqModel) (*users.CompleteSignUpResponse, error)
	CreateCast(userId string, payload CastPayload) error
	NeynarAccountExists(userId string) (bool, error)
}
