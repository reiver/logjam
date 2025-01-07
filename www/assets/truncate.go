package verboten

import (
	"github.com/reiver/go-erorr"
)

func truncateRequestURI(requestPathPrefix string, requestPath string) (string, error) {

	var length int = len(requestPathPrefix)
	if lenRequestPath := len(requestPath); lenRequestPath < length {
		var nada string
		return nada, erorr.Errorf("%s: cannot truncate request-uri — expected its length to be at-least %d bytes long but was actually %d byts long", logtag, length, lenRequestPath)
	}

	var result string = requestPath[len(requestPathPrefix):]
	return result, nil
}
