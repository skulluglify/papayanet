package kornet

import (
  m "PapayaNet/papaya/koala/mapping"
  "PapayaNet/papaya/koala/pp"
  "encoding/json"
  "encoding/xml"
  "errors"
  "github.com/valyala/fasthttp"
  "net/http"
  "net/url"
  "reflect"
  "strconv"
  "strings"
)

// fallback to use default values

func KRequestGetURL(req *http.Request) url.URL {

  URL := url.URL{
    User:     req.URL.User,
    Scheme:   pp.QStr(req.URL.Scheme, "http"),
    Host:     pp.QStr(req.URL.Host, "localhost"),
    Path:     pp.QStr(req.URL.Path, "/"),
    RawQuery: req.URL.RawQuery,
  }

  return URL
}

func KSafeContentTy(contentTy string) (string, string) {

  charset := "UTF-8"

  tokens := strings.Split(contentTy, ";")

  if len(tokens) > 0 {

    contentTy = strings.Trim(tokens[0], "")

    if len(tokens) > 1 {

      token := strings.Trim(tokens[1], " ")
      if strings.HasPrefix("charset", token) {

        tokens = strings.Split(token, "=")

        if len(tokens) > 0 {

          charset = strings.ToUpper(tokens[1])
        }
      }
    }

    if !AvailableCharsets.Contain(charset) {

      charset = "UTF-8"
    }

  } else {

    contentTy = "application/octet-stream"
  }

  return contentTy, charset
}

func KSafeParsingRequestBody(req *fasthttp.Request) (m.KMapImpl, error) {

  //charset := "UTF-8"
  contentTy := string(req.Header.ContentType())
  contentTy, _ = KSafeContentTy(contentTy) // get content-type only

  res := &map[string]any{}

  switch contentTy {

  case "application/json":

    if err := json.Unmarshal(req.Body(), res); err != nil {

      return nil, err
    }

    break

  case "text/xml", "application/xml", "application/atom+xml":

    if err := xml.Unmarshal(req.Body(), res); err != nil {

      return nil, err
    }

    break

  case "multipart/form-data":

    form, err := req.MultipartForm()

    if err != nil {

      return nil, err
    }

    mm := *res

    for k, v := range form.Value {

      if len(v) == 1 {

        mm[k] = v[0]
        continue
      }

      mm[k] = v
    }
  }

  mm := m.KMap(*res)
  return &mm, nil
}

func KSafeParsingBoolean(v any) (bool, error) {

  val := pp.KIndirectValueOf(v)

  if val.IsValid() {

    ty := val.Type()

    switch ty.Kind() {
    case reflect.Bool:

      return val.Bool(), nil

    case reflect.String:

      // try parsing into boolean
      y, err := strconv.ParseBool(val.String())

      if err != nil {

        return false, err
      }

      return y, nil
    }
  }

  return false, errors.New("invalid boolean")
}

func KSafeParsingNumber(v any) (any, error) {

  val := pp.KIndirectValueOf(v)

  if val.IsValid() {

    ty := val.Type()

    switch ty.Kind() {
    case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:

      return val.Int(), nil

    case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:

      return val.Uint(), nil

    case reflect.Float32, reflect.Float64:

      return val.Float(), nil

    case reflect.Complex64, reflect.Complex128:

      return val.Complex(), nil

    case reflect.String:

      k := val.String()

      var n any

      // how about complex ?
      // don't have any idea to parsing with complex

      n, err := strconv.ParseFloat(k, 10)

      if err != nil {

        n, err = strconv.ParseUint(k, 10, 64)

        if err != nil {

          return 0, err
        }
      }

      return n, nil
    }

  }

  return 0, errors.New("invalid number")
}

func KSafeSimpleHeaders(headers url.Values) m.KMapImpl {

  data := map[string]any{}

  for k, v := range headers {

    if len(v) > 0 {

      data[k] = v[0]
    }
  }

  mm := m.KMap(data)
  return &mm
}
