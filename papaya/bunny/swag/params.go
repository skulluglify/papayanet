package swag

import (
  "PapayaNet/papaya/koala"
  m "PapayaNet/papaya/koala/mapping"
  "PapayaNet/papaya/koala/pp"
  "reflect"
)

// once set char be known required or not

func SwagParamPathValid(h string) (bool, string) {

  n := len(h)

  // use dollar sign
  if koala.KStrHasPrefixChar(h, "$") {

    return true, h[1:]
  }

  if koala.KStrHasSuffixChar(h, "$") {

    return true, h[:n-1]
  }

  // use hashtag
  if koala.KStrHasPrefixChar(h, "#") {

    return true, h[1:]
  }

  if koala.KStrHasSuffixChar(h, "#") {

    return true, h[:n-1]
  }

  // is not path
  return false, h
}

func SwagParamRequired(h string) (bool, string) {

  n := len(h)

  if koala.KStrHasPrefixChar(h, "?") {

    return false, h[1:]
  }

  if koala.KStrHasSuffixChar(h, "?") {

    return false, h[:n-1]
  }

  if koala.KStrHasPrefixChar(h, "!") {

    return true, h[1:]
  }

  if koala.KStrHasSuffixChar(h, "!") {

    return true, h[:n-1]
  }

  // default, required
  return true, h
}

func SwagParamsFormatter(mapping any) []m.KMapImpl {

  res := make([]m.KMapImpl, 0)
  val := pp.KIndirectValueOf(mapping)

  if val.IsValid() {

    ty := val.Type()

    switch ty.Kind() {
    case reflect.Map:

      if ty == reflect.TypeOf(m.KMap{}) {

        sample := val.Interface()
        if mm := m.KMapCast(sample); mm != nil {

          for _, enum := range mm.Enums() {

            k, v := enum.Tuple()

            required, token := SwagHeaderRequired(k)
            isPath, name := SwagParamPathValid(token)

            var header m.KMapImpl

            schema := SwagContentFormatter(v)

            header = &m.KMap{
              "in":       pp.LStr(isPath, "path", "query"),
              "name":     name,
              "required": required,
              "schema":   schema,
            }

            // redoc requirment, specific compatible
            header.Put("type", schema.Get("type"))

            res = append(res, header)
          }
        }
      }

      break
    }
  }

  return res
}
