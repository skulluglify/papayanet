package mapping

import (
  "PapayaNet/papaya/koala/pp"
  "reflect"
)

// try parsing into traditional typing, like boolean, int, uint, float, complex, string

func KValueToBool(value any) bool {

  val := pp.KIndirectValueOf(value)

  if val.IsValid() {

    ty := val.Type()

    switch ty.Kind() {
    case reflect.Bool:

      return val.Bool()
    }
  }

  return false
}

func KValueToInt(value any) int64 {

  val := pp.KIndirectValueOf(value)

  if val.IsValid() {

    ty := val.Type()

    switch ty.Kind() {
    case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:

      return val.Int()
    }
  }

  return 0
}

func KValueToUint(value any) uint64 {

  val := pp.KIndirectValueOf(value)

  if val.IsValid() {

    ty := val.Type()

    switch ty.Kind() {
    case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:

      return val.Uint()
    }
  }

  return 0
}

func KValueToFloat(value any) float64 {

  val := pp.KIndirectValueOf(value)

  if val.IsValid() {

    ty := val.Type()

    switch ty.Kind() {
    case reflect.Float32, reflect.Float64:

      return val.Float()
    }
  }

  return 0
}

func KValueToComplex(value any) complex128 {

  val := pp.KIndirectValueOf(value)

  if val.IsValid() {

    ty := val.Type()

    switch ty.Kind() {
    case reflect.Complex64, reflect.Complex128:

      return val.Complex()
    }
  }

  return 0
}

func KValueToString(value any) string {

  val := pp.KIndirectValueOf(value)

  if val.IsValid() {

    ty := val.Type()

    switch ty.Kind() {
    case reflect.String:

      return val.String()
    }
  }

  return ""
}
