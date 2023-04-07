package kornet

import (
  "PapayaNet/papaya/koala/pp"
  "net/http"
  "net/url"
)

// fallback to use default values

func KRequestGetURL(req *http.Request) url.URL {

  URL := url.URL{
    User:     req.URL.User,
    Scheme:   pp.KCOStr(req.URL.Scheme, "http"),
    Host:     pp.KCOStr(req.URL.Host, "localhost"),
    Path:     pp.KCOStr(req.URL.Path, "/"),
    RawQuery: req.URL.RawQuery,
  }

  return URL
}
