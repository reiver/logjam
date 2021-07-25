package statichandler

import (
	"net/http"
)

var Handler http.Handler = http.FileServer(http.Dir("./public/statics/files"))
