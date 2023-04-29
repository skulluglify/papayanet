package cors

import (
  "github.com/valyala/fasthttp"
  "net/http"
  "net/url"
  "skfw/papaya/koala/mapping"
  "skfw/papaya/koala/pp"
  "strings"
)

// ref: https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS

type Consumer struct {
  URL         *url.URL
  Methods     []string // GET, POST, PUT, ...
  Headers     []string
  Credentials bool
  MaxAge      int // delta - seconds, ex: 86400 > 24 hours
}

type ConsumerImpl interface {
  Init(URL *url.URL, methods []string, headers []string, credentials bool, maxAge int) error
  Header() *http.Header
  AcceptMethod(method string) bool
  CheckOrigin(origin string) bool
  Origin(origin string) ConsumerImpl
  Check(method string, origin string) bool
  Apply(res *fasthttp.Response)
}

func ConsumerNew(URL *url.URL, methods []string, headers []string, credentials bool, maxAge int) (ConsumerImpl, error) {

  consumer := &Consumer{}
  if err := consumer.Init(URL, methods, headers, credentials, maxAge); err != nil {

    return nil, err
  }

  return consumer, nil
}

func (c *Consumer) Init(URL *url.URL, methods []string, headers []string, credentials bool, maxAge int) error {

  // normalize url by remove prefix of www.
  URL.Host, _ = strings.CutPrefix(URL.Host, "www.")

  c.URL = URL
  c.Methods = methods
  c.Headers = headers
  c.Credentials = credentials
  c.MaxAge = maxAge

  return nil
}

func (c *Consumer) Header() *http.Header {

  // Access-Control-Request-Method: POST
  // Access-Control-Request-Headers: Content-Type
  // Accept: text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8
  // Accept-Language: en-us,en;q=0.5
  // Accept-Encoding: gzip,deflate
  // Accept-Charset: ISO-8859-1,utf-8;q=0.7,*;q=0.7
  // Connection: keep-alive

  // Access-Control-Allow-Origin
  // Access-Control-Allow-Methods
  // Access-Control-Allow-Headers
  // Access-Control-Allow-Credentials
  // Access-Control-Max-Age
  // Vary: Accept-Encoding, Origin
  // Content-Encoding: gzip
  // Keep-Alive: timeout=2, max=100
  // Connection: Keep-Alive

  header := &http.Header{}

  // fixing a problem if replace current origin
  // got a problem if origin a have prefix of www.
  header.Add("Access-Control-Allow-Origin", SafeURL(c.URL))

  header.Add("Access-Control-Allow-Methods", SafeMethods(c.Methods))
  header.Add("Access-Control-Allow-Headers", SafeHeaders(c.Headers))
  header.Add("Access-Control-Allow-Credentials", SafeCredentials(c.Credentials))
  header.Add("Access-Control-Max-Age", SafeMaxAge(c.MaxAge))
  // header.Add("Vary", "Accept-Encoding, Origin") // accept encoding for future

  return header
}

func (c *Consumer) AcceptMethod(method string) bool {

  if !CheckMethodAvail(method) {

    return false
  }

  if !mapping.Keys(c.Methods).Contain(method) {

    c.Methods = append(c.Methods, method) // added new method
  }

  return true // method has been added
}

func (c *Consumer) CheckOrigin(origin string) bool {

  URL, err := url.Parse(origin)
  if err != nil {

    return false // can't parse origin
  }

  // source scheme fallback into credential use case
  if URL.Scheme == pp.QStr(c.URL.Scheme, "https") {

    // normalize url by remove prefix of www.
    URL.Host, _ = strings.CutPrefix(URL.Host, "www.")

    // must same as origin, safe compare www.
    if URL.Host == c.URL.Host {

      return true
    }
  }

  return false // fallback into default value
}

func (c *Consumer) Origin(origin string) ConsumerImpl {

  URL, err := url.Parse(origin)
  if err != nil {

    return nil // bypass, unable to parse url
  }

  // copy data
  return &Consumer{
    URL:         URL,
    Methods:     c.Methods,
    Headers:     c.Headers,
    Credentials: c.Credentials,
    MaxAge:      c.MaxAge,
  }
}

func (c *Consumer) Check(method string, origin string) bool {

  // check method is granted or denied
  if !mapping.Keys(c.Methods).Contain(method) {

    return false
  }

  // check origin
  return c.CheckOrigin(origin)
}

func (c *Consumer) Apply(res *fasthttp.Response) {

  var header *http.Header
  var value string

  // get header from consumer information
  header = c.Header()

  for _, key := range Headers {

    if value = header.Get(key); value != "" {

      res.Header.Set(key, value)
    }
  }
}
