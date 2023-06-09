package swag

import (
  "reflect"
  "skfw/papaya/koala"
  m "skfw/papaya/koala/mapping"
  "skfw/papaya/koala/pp"
)

// once set char be known required or not

func SwagHeaderRequired(h string) (bool, string) {

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

func SwagHeadersFormatter(mapping any) []m.KMapImpl {

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

            required, name := SwagHeaderRequired(k)

            var header m.KMapImpl

            schema := SwagContentFormatter(v)

            header = &m.KMap{
              "in":       "header",
              "name":     name,
              "required": required,
              "schema":   schema,
              "type":     "object",
            }

            // redoc requirement, specific compatible
            //header.Put("type", schema.Get("type"))

            res = append(res, header)
          }
        }
      }

      break
    }
  }

  return res
}
