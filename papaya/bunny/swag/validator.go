package swag

import (
  "encoding/json"
  "errors"
  "fmt"
  "net/url"
  "reflect"
  "skfw/papaya/koala/kio/leaf"
  "skfw/papaya/koala/kornet"
  m "skfw/papaya/koala/mapping"
  "skfw/papaya/koala/pp"
  "strconv"

  "github.com/gofiber/fiber/v2"
)

// request validation, only on type json, xml, form

type SwagRequestValidator struct {
  *fiber.Ctx
  exp m.KMapImpl
}

type SwagRequestValidatorImpl interface {
  Init(exp m.KMapImpl, ctx *fiber.Ctx)
  Validation() (*kornet.Request, error)
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

func (v *SwagRequestValidator) Validation() (*kornet.Request, error) {

  // try getting content-type and charset

  request := &kornet.Request{}

  //charset := "UTF-8"
  bCTy := string(v.Ctx.Request().Header.ContentType()) // body request content-type
  bCTy, _ = kornet.KSafeContentTy(bCTy)

  // -- end

  var body m.KMapImpl
  body = &m.KMap{}

  req := v.Ctx.Request()

  //params := m.KMapCast(v.exp.Get("params"))
  //headers := m.KMapCast(v.exp.Get("headers"))

  // validation the request body
  if content := m.KMapCast(v.exp.Get("request.body")); content != nil {

    cTys := content.Keys()

    if cTys.Contain(bCTy) {

      mm, err := kornet.KSafeParsingRequestBody(req)

      if err != nil {

        return request, fmt.Errorf("data format does not match mime type %s", bCTy)
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
      bCTy = "application/json"

      mm := m.KMap(*data)
      body = &mm
    }

    var RequestBodySameAsContentTypeRequired bool
    var RequestBodyKeySameAsSampleKey bool

    for _, cTy := range cTys {

      if cTy != bCTy {

        continue
      }

      schema := content.Get(cTy + ".schema") // ex: application/json.schema

      if schema == nil { // schema from expectation

        // {{content-type}}.schema not found
        // not set up schema for handle it

        return request, errors.New("schema from request body is not implemented")

      }

      if mm := m.KMapCast(schema); mm != nil {

        for _, enum := range mm.Tree().Enums() { // schema required

          schemaKey, schemaType := enum.Tuple()

          rt := SwagUniversalReType(schemaKey)

          compareSampleKey, err := CompareSampleKeyNew(schemaKey)

          if err != nil {

            return request, errors.New("unable to compare with sample key from schema object")
          }

          RequestBodyKeySameAsSampleKey = false

          for _, enum := range body.Tree().Enums() { // schema request body

            bK, bV := enum.Tuple()

            if !compareSampleKey.Check(bK) {

              continue
            }

            RequestBodyKeySameAsSampleKey = true

            switch rt {

            case "bool", "boolean":

              y, err := kornet.KSafeParsingBoolean(bV)

              if err != nil {

                return request, fmt.Errorf("key `%s` either not set or is not a boolean in request body", schemaKey)
              }

              // update value from request body
              body.Set(bK, y)

              break

            case "int", "number", "integer", "byte": // byte -> uint8

              // try parsing if not number
              n, err := kornet.KSafeParsingNumber(bV)

              if err != nil {

                return request, fmt.Errorf("key `%s` either not set or is not a number in request body", schemaKey)
              }

              // update value from request body
              body.Set(bK, n)

              break

            case "str", "text", "string":

              // check is string or not
              if tx := m.KValueToString(bV); tx == "" {

                return request, fmt.Errorf("key `%s` either not set or is not a string in request body", schemaKey)
              }

              break

            case "array":

              schemaTypeOf := reflect.TypeOf(schemaType) // type of type

              // test body value
              val := pp.KIndirectValueOf(bV)

              if val.IsValid() {

                bodyTypeOf := val.Type() // get type of body value

                schemaTypeOfKey := schemaTypeOf.Elem() // type of type element

                // get a element type in universal type
                uniTypeSchemaTypeOfKey := SwagUniversalReType(schemaTypeOfKey.Name()) // re type, type of type

                switch bodyTypeOf.Kind() {
                case reflect.Array, reflect.Slice: // check body value is array, or slice

                  for i := 0; i < val.Len(); i++ {

                    indexBodyValue := pp.KIndirectValueOf(val.Index(i))

                    if indexBodyValue.IsValid() {

                      typeIndexBodyValue := indexBodyValue.Type()

                      // that problem is,
                      // tte = map[k]v
                      // vtt = map[k]v
                      // tte != vtt
                      // same as a array or slice too

                      // fast compare
                      if schemaTypeOfKey.Kind() != typeIndexBodyValue.Kind() {

                        return request, fmt.Errorf("key `%s` either not set or is not a array<%s> in request body", schemaKey, uniTypeSchemaTypeOfKey)
                      }
                    }
                  }

                default:

                  return request, fmt.Errorf("key `%s` either not set or is not a array<%s> in request body", schemaKey, uniTypeSchemaTypeOfKey)
                }
              }

              break

            case "object":

              // skip map, not null

              if bV == nil {

                return request, fmt.Errorf("key `%s` is null in request body", schemaKey)
              }

              break
            }
          }

          if !RequestBodyKeySameAsSampleKey {

            return request, fmt.Errorf("key `%s` is null in request body", schemaKey)
          }
        }
      }

      RequestBodySameAsContentTypeRequired = true
      break
    }

    if !RequestBodySameAsContentTypeRequired {

      return request, fmt.Errorf("content-type from request body is not supported")
    }
  }

  // re-packing into json format
  request.Body = leaf.KMakeBuffer([]byte(body.JSON()))

  header := map[string]any{}

  if content := m.KMapCast(v.exp.Get("request.headers")); content != nil {

    h := v.Ctx.GetReqHeaders()

    for _, enum := range content.Enums() {

      k, t := enum.Tuple()

      required, token := SwagHeaderRequired(k)

      if required {

        if p, ok := h[token]; ok {

          switch t {
          case "bool", "boolean":

            y, err := strconv.ParseBool(p)

            if err != nil {

              return request, fmt.Errorf("key `%s` either not set or is not a boolean in header", k)
            }

            header[token] = y

            break

          case "int", "number", "integer", "byte":

            // try parsing if not number
            n, err := kornet.KSafeParsingNumber(p)

            if err != nil {

              return request, fmt.Errorf("key `%s` either not set or is not a number in request header", k)
            }

            header[token] = n

            break

          case "str", "text", "string":

            header[token] = p

            break
          }
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

  mp := m.KMap(paths)
  request.Path = &mp

  request.Query = queries

  return request, nil
}
