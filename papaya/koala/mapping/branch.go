package mapping

import (
  "fmt"
  "reflect"
  "skfw/papaya/koala/pp"
  "strings"
)

// TODO: fixing new future

func Branch(name string, mapping any) any {

  val := pp.KIndirectValueOf(mapping)
  tokens := strings.Split(name, ".")

  if val.IsValid() {

    ty := val.Type()

    switch ty.Kind() {
    case reflect.Map:

      if ty.Key().Kind() == reflect.String {

        for _, token := range tokens {

          fmt.Println(token)
          panic("not implemented yet")
        }
      }
    }
  }

  return nil
}
