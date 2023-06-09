package cors

import (
  "errors"
  "github.com/gofiber/fiber/v2"
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
  Header(origin string, requestMethods []string, requestHeaders []string) (*http.Header, error)
  AcceptMethod(method string) bool
  Check(method string, origin string) bool
  Apply(ctx *fiber.Ctx) *fiber.Ctx
}

func ConsumerNew(URL *url.URL, methods []string, headers []string, credentials bool, maxAge int) (ConsumerImpl, error) {

  consumer := &Consumer{}
  if err := consumer.Init(URL, methods, headers, credentials, maxAge); err != nil {

    return nil, err
  }

  return consumer, nil
}

func (c *Consumer) Init(URL *url.URL, methods []string, headers []string, credentials bool, maxAge int) error {

  if URL != nil {

    // normalize url by remove prefix of www.
    URL.Host, _ = strings.CutPrefix(URL.Host, "www.")
  }

  c.URL = URL
  c.Methods = methods
  c.Headers = headers
  c.Credentials = credentials
  c.MaxAge = maxAge

  return nil
}

func (c *Consumer) Header(origin string, requestMethods []string, requestHeaders []string) (*http.Header, error) {

  var err error

  var URL *url.URL
  var Origin string

  var header *http.Header

  var methods []string
  var headers []string

  header = &http.Header{}

  // try fallback with current URL
  if origin != "" && origin != "*" {

    // check origin is URL
    URL, err = url.Parse(origin)
    if err != nil {

      return header, errors.New("unable to parse URL from origin")
    }

    if URL.Scheme != "" && URL.Host != "" {

      Origin = URL.Scheme + "://" + URL.Host

    } else {

      return header, errors.New("undefined scheme or host from URL")
    }

  } else {

    if c.URL != nil {

      // fixing a problem if replace current origin
      // got a problem if origin a have prefix of www.
      Origin = SafeURL(c.URL)

    } else {

      // origin asterisk allowed
      Origin = "*"
    }
  }

  // normalize methods with request methods
  methods = NormListBySources(c.Methods, requestMethods)

  // normalize headers with request headers
  headers = NormListBySources(c.Headers, requestHeaders)

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

  // asterisk

  header.Add("Access-Control-Allow-Origin", Origin)
  header.Add("Access-Control-Allow-Methods", SafeMethods(methods))
  header.Add("Access-Control-Allow-Headers", SafeHeaders(headers))
  header.Add("Access-Control-Allow-Credentials", SafeCredentials(c.Credentials))
  header.Add("Access-Control-Max-Age", SafeMaxAge(c.MaxAge))
  // header.Add("Vary", "Accept-Encoding, Origin") // accept encoding for future

  return header, nil
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

func (c *Consumer) Check(method string, origin string) bool {

  // asterisk

  // passing method by empty string
  if method != "" && method != "*" {

    // check method is granted or denied
    if !mapping.Keys(c.Methods).Contain(method) {

      return false
    }
  }

  if c.URL != nil {

    if origin != "" && origin != "*" {

      // check origin
      URL, err := url.Parse(origin)
      if err != nil {

        return false // can't parse origin
      }

      // source scheme fallback into credential use case
      if URL.Scheme == pp.Qstr(c.URL.Scheme, "https") {

        // normalize url by remove prefix of www.
        URL.Host, _ = strings.CutPrefix(URL.Host, "www.")

        // must same as origin, safe compare www.
        if URL.Host == c.URL.Host {

          return true
        }
      }
    }

  } else {

    // enable all if current URL is NULL
    return true
  }

  return false // fallback into default value
}

func (c *Consumer) Apply(ctx *fiber.Ctx) *fiber.Ctx {

  // catch request and response method
  req, res := ctx.Request(), ctx.Response()

  // must be checked first

  var value string
  var header *http.Header

  var Origin string
  var Referer string
  var Method string

  var methods []string
  var headers []string

  // may fix with, curr method and req method
  methods = make([]string, 0)

  Origin = string(req.Header.Peek("Origin"))
  Referer = string(req.Header.Peek("Referer"))

  if Referer != "" {

    // fix issue about strict cross-origin
    Referer, _ = RefererIntoOrigin(Referer)
    Origin = pp.Qstr(Origin, Referer)
  }

  // :|
  //method = "*"

  // noop, don't have any idea for used

  Method = string(req.Header.Peek("Access-Control-Request-Method"))

  // set current method and request method
  methods = append(methods, strings.ToUpper(string(req.Header.Method())))

  if Method != "" {

    methods = append(methods, Method)
  }

  // :|
  headers = strings.Split(string(req.Header.Peek("Access-Control-Request-Headers")), ",")

  // passing error
  // get header from consumer information
  header, _ = c.Header(Origin, methods, headers)

  for _, key := range Headers {

    if value = header.Get(key); value != "" {

      res.Header.Set(key, value)
    }
  }

  return ctx
}
