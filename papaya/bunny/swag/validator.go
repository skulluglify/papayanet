package swag

import (
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

  var err error

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

    ////////// try parsing //////////

    cTys := content.Keys()

    body, err = kornet.KSafeParsingRequestBody(req)

    if err != nil {

      return request, fmt.Errorf("data format does not match mime type %s", bCTy)
    }

    // fake result as JSON

    bCTy = "application/json"

    ////////// try parsing //////////

    var requestBodySameAsContentTypeRequired bool

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

        var found bool
        var required bool
        var compareSampleKey CompareSampleKeyImpl

        for _, schemaEnum := range mm.Tree().Enums() { // schema required

          schemaKey, schemaType := schemaEnum.Tuple()

          compareSampleKey, err = CompareSampleKeyNew(schemaKey) // watch requireable keys
          required, schemaKey = SwagRequired(schemaKey)          // remove "?" "!" char

          rt := SwagUniversalReType(schemaType)

          if err != nil {

            return request, errors.New("unable to compare with sample key from schema object")
          }

          found = false

          for _, enum := range body.Tree().Enums() { // schema request body

            key, bV := enum.Tuple()

            if !compareSampleKey.ReCheck(key) {

              continue
            }

            found = true

            switch rt {

            case "bool", "boolean":

              y, err := kornet.KSafeParsingBoolean(bV)

              if required {

                if err != nil {

                  return request, fmt.Errorf("key `%s` either not set or is not a boolean in request body", schemaKey)
                }
              }

              // update value from request body
              body.Set(key, y)

              break

            case "int", "number", "numeric", "integer", "byte": // byte -> uint8

              // try parsing if not number
              n, err := kornet.KSafeParsingNumber(bV)

              if required {

                if err != nil {

                  return request, fmt.Errorf("key `%s` either not set or is not a number in request body", schemaKey)
                }
              }

              // update value from request body
              body.Set(key, n)

              break

            case "str", "string", "text", "word":

              tx := m.KValueToString(bV)

              if required {

                if tx == "" {

                  return request, fmt.Errorf("key `%s` either not set or is not a string in request body", schemaKey)
                }
              }

              body.Set(key, tx)

              break

            case "array":

              schemaTypeOf := reflect.TypeOf(schemaType) // type of type

              // test body value
              val := pp.KIndirectValueOf(bV)

              if required {

                if val.IsValid() {

                  bodyTypeOf := val.Type() // get a type of body value

                  schemaTypeOfKey := schemaTypeOf.Elem() // type of type element

                  var uniTypeSchemaTypeOfKey string
                  // get a element type in universal type
                  switch schemaTypeOfKey.Kind() {
                  case reflect.String:
                    uniTypeSchemaTypeOfKey = SwagUniversalReType(schemaTypeOfKey.Name()) // re type, type of type
                  case reflect.Struct, reflect.Map:
                    uniTypeSchemaTypeOfKey = "object"
                  default:
                    uniTypeSchemaTypeOfKey = "null"
                  }

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
                        // same as an array or slice too
                        // fast compare

                        // check similarity kind of type like object
                        check := func() bool {

                          // map, or struct like object
                          switch schemaTypeOfKey.Kind() {
                          case reflect.Struct, reflect.Map:

                            // map, or struct like object
                            switch typeIndexBodyValue.Kind() {
                            case reflect.Struct, reflect.Map:
                              return true
                            }

                          default:

                            // compare kind
                            return schemaTypeOfKey.Kind() == typeIndexBodyValue.Kind()
                          }

                          return false
                        }()

                        if !check {

                          return request, fmt.Errorf("key `%s` either not set or is not a array<%s> in request body", schemaKey, uniTypeSchemaTypeOfKey)
                        }
                      }
                    }

                  default:

                    return request, fmt.Errorf("key `%s` either not set or is not a array<%s> in request body", schemaKey, uniTypeSchemaTypeOfKey)
                  }
                }
              }

              break

            case "object":

              // skip map, not null

              if required {

                if bV == nil {

                  return request, fmt.Errorf("key `%s` is null in request body", schemaKey)
                }
              }

              break
            }
          }

          if required {

            if !found {

              return request, fmt.Errorf("key `%s` is null in request body", schemaKey)
            }
          }
        }
      }

      requestBodySameAsContentTypeRequired = true
      break
    }

    if !requestBodySameAsContentTypeRequired {

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

      if p, ok := h[token]; ok {

        switch t {
        case "bool", "boolean":

          y, err := strconv.ParseBool(p)

          if required {

            if err != nil {

              return request, fmt.Errorf("key `%s` either not set or is not a boolean in header", k)
            }
          }

          header[token] = y

          break

        case "int", "number", "integer", "byte":

          // try parsing if not number
          n, err := kornet.KSafeParsingNumber(p)

          if required {

            if err != nil {

              return request, fmt.Errorf("key `%s` either not set or is not a number in request header", k)
            }
          }

          header[token] = n

          break

        case "str", "string", "text", "word":

          if required {

            if p == "" {

              return request, fmt.Errorf("key `%s` not found in header", k)
            }
          }

          header[token] = p

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

    var found bool

    for _, enum := range content.Enums() {

      k, t := enum.Tuple()

      //fmt.Println(k, t)

      // params

      required, token := SwagParamRequired(k)
      isPath, name := SwagParamPathValid(token)

      found = false

      if isPath {

        // path

        for p, q := range params {

          if p == name {

            // try parsing into boolean, number

            switch t {
            case "bool", "boolean":

              y, err := strconv.ParseBool(q)

              if required {

                if err != nil {

                  return request, fmt.Errorf("key `%s` either not set or is not a boolean in path", name)
                }
              }

              paths[name] = y
              found = true

              break

            case "int", "number", "integer", "byte":

              // try parsing if not number
              n, err := kornet.KSafeParsingNumber(q)

              if required {

                if err != nil {

                  return request, fmt.Errorf("key `%s` either not set or is not a number in path", name)
                }
              }

              paths[name] = n
              found = true

              break

            case "str", "string", "text", "word":

              if required {

                if q == "" {

                  return request, fmt.Errorf("key `%s` not found in path", name)
                }
              }

              paths[name] = q
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

        switch t {
        case "bool", "boolean":

          y, err := strconv.ParseBool(q)

          if required {

            if err != nil {

              return request, fmt.Errorf("key `%s` either not set or is not a boolean in query", name)
            }
          }

          queries.Set(name, y)
          found = true

          break

        case "int", "number", "integer", "byte":

          // try parsing if not number
          n, err := kornet.KSafeParsingNumber(q)

          if required {

            if err != nil {

              return request, fmt.Errorf("key `%s` either not set or is not a number in query", k)
            }
          }

          queries.Set(name, n)
          found = true

          break

        case "str", "string", "text", "word":

          if required {

            if q == "" {

              return request, fmt.Errorf("key `%s` not found in query", name)
            }
          }

          queries.Set(name, q)
          found = true

          break
        }

        if !found {

          queries.Set(name, q)
        }
      }
    }
  }

  mp := m.KMap(paths)
  request.Path = &mp

  request.Query = queries

  return request, nil
}
