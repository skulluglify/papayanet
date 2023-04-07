package kornet

import (
  "PapayaNet/papaya/koala/io/leaf"
  "net/http"
  "net/url"
)

type KHeaders http.Header
type KParams url.Values

type KRequest struct {
  Method  string
  URL     url.URL
  Headers KHeaders
  Params  KParams
  Body    leaf.KBufferImpl
}

type KResponse struct {
  Headers KHeaders
  Params  KParams
  Body    leaf.KBufferImpl
}
