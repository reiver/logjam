package logjamlink

import (
	"fmt"
)

func LogJamLink(host string, name string) string {
	if "" == host {
		return ""
	}
	if "" == name {
		return ""
	}

	return fmt.Sprintf("logjam://%s@%s/conf", name, host)
}
