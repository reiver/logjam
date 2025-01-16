package db

import (
	"net/url"

	"github.com/reiver/go-path"
)

func createURL(baseURL string, apipath string) string {
	apipath = path.Canonical(apipath)
	if "." == apipath {
		apipath = ""
	}
	if 0 < len(apipath) && '/' == apipath[0] {
		apipath = apipath[1:]
	}

	urloc, err := url.Parse(baseURL)
	if nil != err {
		return ""
	}
	if "" == urloc.Scheme {
		return ""
	}
	if "" == urloc.Host {
		return ""
	}

	var result string
	{
		var buffer [256]byte
		var p []byte = buffer[0:0]

		p = append(p, urloc.Scheme...)
		p = append(p, "://"...)
		if nil != urloc.User {
			p = append(p, urloc.User.String()...)
			p = append(p, '@')
		}
		p = append(p, urloc.Host...)

		var urlpath string = path.Join(urloc.Path, apipath)
		if "." == urlpath {
			urlpath = "/"
		}
		if 0 < len(urlpath) && '/' != urlpath[0] {
			urlpath = "/" + urlpath
		}
		p = append(p, urlpath...)

		result = string(p)
	}

	return result
}
