package swag

import (
  "PapayaNet/papaya/koala/kio/leaf"
  "PapayaNet/papaya/koala/kornet"
  m "PapayaNet/papaya/koala/mapping"
  "PapayaNet/papaya/koala/pp"
  "encoding/json"
  "errors"
  "fmt"
  "github.com/gofiber/fiber/v2"
  "net/url"
  "reflect"
  "strconv"
)

// request validation, only on type json, xml, form

type SwagRequestValidator struct {
  *fiber.Ctx
  exp m.KMapImpl
}

type SwagRequestValidatorImpl interface {
  Init(exp m.KMapImpl, ctx *fiber.Ctx)
  Validation() (kornet.Request, error)
}

func SwagRequestValidatorNew(exp m.KMapImpl, ctx *fiber.Ctx) SwagRequestValidatorImpl {

  v := &SwagRequestValidator{}
  v.Init(exp, ctx)
  return v
}

func (v *SwagRequestValidator) Init(exp m.KMapImpl, ctx *fiber.Ctx) {

  v.Ctx = ctx
  v.exp = exp
}

func (v *SwagRequestValidator) Validation() (kornet.Request, error) {

  // try getting content-type and charset

  request := kornet.Request{}

  //charset := "UTF-8"
  contentTy := string(v.Ctx.Request().Header.ContentType())
  contentTy, _ = kornet.KSafeContentTy(contentTy)

  // -- end

  var body m.KMapImpl
  body = &m.KMap{}

  req := v.Ctx.Request()

  //params := m.KMapCast(v.exp.Get("params"))
  //headers := m.KMapCast(v.exp.Get("headers"))

  // validation the request body
  if content := m.KMapCast(v.exp.Get("request.body")); content != nil {

    if content.Keys().Contain(contentTy) {

      mm, err := kornet.KSafeParsingRequestBody(req)

      if err != nil {

        return request, fmt.Errorf("data format does not match mime type %s", contentTy)
      }

      body = mm

    } else {

      // use lucky
      // try parsing with format json
      data := &map[string]any{}
      if err := json.Unmarshal(req.Body(), data); err != nil {

        return request, errors.New("failed parsing into data json")
      }

      // update current content-type
      contentTy = "application/json"

      mm := m.KMap(*data)
      body = &mm
    }

    contentTys := content.Keys()

    var found bool
    for _, cTy := range contentTys {

      if cTy == contentTy {

        schema := content.Get(cTy + ".schema")

        if mm := m.KMapCast(schema); mm != nil {

          // tree iteration
          iter := mm.Tree().Iterable()
          // iteration all values, bool, number, string
          for next := iter.Next(); next.HasNext(); next = next.Next() {

            enum := next.Enum()
            k, t := enum.Tuple()

            rt := SwagUniversalReType(t)

            p := body.Get(k)

            // foolish reason if value set into null

            if p != nil {

              switch rt {

              case "bool", "boolean":

                y, err := kornet.KSafeParsingBoolean(p)

                if err != nil {

                  return request, fmt.Errorf("key `%s` either not set or is not a boolean in request body", k)
                }

                body.Set(k, y)

                break

              case "int", "number", "integer", "byte": // byte -> uint8

                // try parsing if not number
                n, err := kornet.KSafeParsingNumber(p)

                if err != nil {

                  return request, fmt.Errorf("key `%s` either not set or is not a number in request body", k)
                }

                // maybe origin is string
                body.Set(k, n)

                break

              case "str", "text", "string":

                // check is string or not
                if tx := m.KValueToString(p); tx == "" {

                  return request, fmt.Errorf("key `%s` either not set or is not a string in request body", k)
                }

                break

              case "array":

                tt := reflect.TypeOf(t)
                val := pp.KIndirectValueOf(p)

                if val.IsValid() {

                  ty := val.Type()

                  tte := tt.Elem()

                  // get a element type in universal type
                  rTte := SwagUniversalReType(tte.Name())

                  switch ty.Kind() {
                  case reflect.Array, reflect.Slice:

                    for i := 0; i < val.Len(); i++ {

                      vt := pp.KIndirectValueOf(val.Index(i))

                      if vt.IsValid() {

                        vtt := vt.Type()

                        // that problem is,
                        // tte = map[k]v
                        // vtt = map[k]v
                        // tte != vtt
                        // same as a array or slice too

                        // fast compare
                        if tte.Kind() != vtt.Kind() {

                          return request, fmt.Errorf("key `%s` either not set or is not a array<%s> in request body", k, rTte)
                        }
                      }
                    }

                  default:

                    return request, fmt.Errorf("key `%s` either not set or is not a array<%s> in request body", k, rTte)
                  }
                }

                break

              case "object":

                break
              }

              continue
            }

            return request, fmt.Errorf("key `%s` is null in request body", k)

          }
        }

        found = true
        break
      }
    }

    if !found {

      return request, errors.New("mime type is not registered")
    }
  }

  // re-packing into json format
  request.Body = leaf.KMakeBuffer([]byte(body.JSON()))

  var header map[string]any

  if content := m.KMapCast(v.exp.Get("request.headers")); content != nil {

    h := v.Ctx.GetReqHeaders()

    for _, enum := range content.Enums() {

      k, t := enum.Tuple()

      if p, ok := h[k]; ok {

        switch t {
        case "bool", "boolean":

          y, err := strconv.ParseBool(p)

          if err != nil {

            return request, fmt.Errorf("key `%s` either not set or is not a boolean in header", k)
          }

          header[k] = y

          break

        case "int", "number", "integer", "byte":

          // try parsing if not number
          n, err := kornet.KSafeParsingNumber(p)

          if err != nil {

            return request, fmt.Errorf("key `%s` either not set or is not a number in request header", k)
          }

          header[k] = n

          break

        case "str", "text", "string":

          header[k] = p

          break
        }
      }
    }
  }

  mm := m.KMap(header)
  request.Header = &mm

  paths := map[string]any{}
  queries := kornet.Query{}

  if content := m.KMapCast(v.exp.Get("request.params")); content != nil {

    params := v.Ctx.AllParams()
    rawQuery := v.Ctx.Context().QueryArgs().String()
    query, _ := url.ParseQuery(rawQuery)

    for _, enum := range content.Enums() {

      k, t := enum.Tuple()

      //fmt.Println(k, t)

      // params

      required, token := SwagParamRequired(k)
      isPath, name := SwagParamPathValid(token)

      if required {

        if isPath {

          // path

          for p, q := range params {

            if p == name {

              // try parsing into boolean, number

              var found bool

              switch t {
              case "bool", "boolean":

                y, err := strconv.ParseBool(q)

                if err != nil {

                  return request, fmt.Errorf("key `%s` either not set or is not a boolean in path", name)
                }

                paths[name] = y
                found = true

                break

              case "int", "number", "integer", "byte":

                // try parsing if not number
                n, err := kornet.KSafeParsingNumber(q)

                if err != nil {

                  return request, fmt.Errorf("key `%s` either not set or is not a number in path", name)
                }

                paths[name] = n
                found = true

                break
              }

              if !found {

                paths[name] = q
              }
            }
          }

        } else {

          // query
          q := query.Get(name)

          var found bool

          switch t {
          case "bool", "boolean":

            y, err := strconv.ParseBool(q)

            if err != nil {

              return request, fmt.Errorf("key `%s` either not set or is not a boolean in query", name)
            }

            queries.Set(name, y)
            found = true

            break

          case "int", "number", "integer", "byte":

            // try parsing if not number
            n, err := kornet.KSafeParsingNumber(q)

            if err != nil {

              return request, fmt.Errorf("key `%s` either not set or is not a number in query", k)
            }

            queries.Set(name, n)
            found = true

            break
          }

          // bypass if a string
          // just parsing check on int, boolean

          if !found {

            queries.Set(name, q)
          }
        }
      }
    }
  }

  pp := m.KMap(paths)
  request.Path = &pp

  request.Query = queries

  return request, nil
}
