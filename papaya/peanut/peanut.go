package peanut

import (
	"net/http"
	"net/url"
)

const (
	PN_PEANUT_METHOD_GET = iota
	PN_PEANUT_METHOD_HEAD
	PN_PEANUT_METHOD_POST
	PN_PEANUT_METHOD_PUT
	PN_PEANUT_METHOD_DELETE
	PN_PEANUT_METHOD_CONNECT
	PN_PEANUT_METHOD_OPTIONS
	PN_PEANUT_METHOD_TRACE
)

func PnGetURLFromRequest(req *http.Request) string {

	scheme := "http"
	host := "127.0.0.1"
	path := "/"
	query := ""

	if req.URL.Scheme != "" {

		scheme = req.URL.Scheme
	}

	if req.URL.Host != "" {

		host = req.URL.Host
	}

	if req.URL.Path != "" {

		path = req.URL.Path
	}

	if req.URL.RawQuery != "" {

		query = req.URL.RawQuery
	}

	URL := url.URL{
		Scheme:   scheme,
		Host:     host,
		Path:     path,
		RawQuery: query,
	}

	return URL.String()
}
