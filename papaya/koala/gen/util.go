package gen

import (
  "reflect"
  "skfw/papaya/koala/pp"
)

func KMapHunt(value any) bool {

  valueOf := pp.KIndirectValueOf(value)

  if valueOf.IsValid() {

    switch valueOf.Kind() {
    case reflect.Array, reflect.Slice, reflect.Map, reflect.Struct:
      return true
    }
  }

  return false
}
