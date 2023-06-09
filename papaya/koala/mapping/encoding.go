package mapping

import (
  "reflect"
  "skfw/papaya/koala"
  "skfw/papaya/koala/pp"
  "strconv"
)

func KMapEncodeJSON(mapping any) string {

  val := pp.KIndirectValueOf(mapping)

  if val.IsValid() {

    ty := val.Type()
    switch ty.Kind() {

    case reflect.Array, reflect.Slice:

      var res string

      n := val.Len()

      for i := 0; i < n; i++ {

        v := val.Index(i).Interface()

        temp := pp.Qstr(koala.KStrRepr(v), KMapEncodeJSON(v))

        if i+1 < n {

          res += temp + ","

        } else {

          res += temp
        }
      }

      return "[" + res + "]"

    case reflect.Map, reflect.Struct:

      var res string

      enums := KMapEnums(val)
      n := len(enums)

      for i, enum := range enums {

        k, v := enum.Tuple()

        temp := pp.Qstr(koala.KStrRepr(v), KMapEncodeJSON(v))

        token := strconv.Quote(k) + ":" + temp

        if i+1 < n {

          res += token + ","

        } else {

          res += token
        }
      }

      return "{" + res + "}"
    }
  }

  return "null"
}
