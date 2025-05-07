package rest

import (
	"github.com/reiver/logjam/lib/users"
	"net/http"
)

func HandleIfErr(rw http.ResponseWriter, err error, status int) bool {
	if err == nil {
		return false
	}

	Error(rw, err, status)
	return true
}

var UserDataCtxKey = "userDataKey"

func GetUser(request *http.Request) *users.UserDTO {
	user, ok := request.Context().Value(UserDataCtxKey).(*users.UserDTO) // whatever your type is
	if !ok {
		return nil
	}
	return user
}
