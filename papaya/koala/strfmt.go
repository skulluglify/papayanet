package koala

import (
  "PapayaNet/papaya/koala/pp"
  "PapayaNet/papaya/panda"
  "reflect"
  "strconv"
)

// Method String to Number

func KStrToNum(value string) int {

  // try convert
  if v, err := strconv.Atoi(value); err == nil {

    return v
  }

  // default value
  return 0
}

// Method convert String to Boolean

func KStrToBool(value string) bool {

  // try convert
  if v, err := strconv.ParseBool(value); err == nil {

    return v
  }

  // default value
  return false
}

// Method convert String into Bytes

func KStrToBytes(value string) []byte {

  return []byte(value)
}

func KStrZeroFill(text string, s int) string {

  var zeros string

  k := panda.Min(len(text), s)
  z := s - k

  for i := 0; i < z; i++ {

    zeros += "0"
  }

  return zeros + text[:k]
}

func KStrPadStart(text string, s int) string {

  var pads string

  k := panda.Min(len(text), s)
  z := s - k

  for i := 0; i < z; i++ {

    pads += " "
  }

  return pads + text[:k]
}

func KStrPadEnd(text string, s int) string {

  var pads string

  k := panda.Min(len(text), s)
  z := s - k

  for i := 0; i < z; i++ {

    pads += " "
  }

  return text[:k] + pads
}

func KStrFmt(value any) string {

  val := pp.KIndirectValueOf(value)

  if val.IsValid() {

    ty := val.Type()

    switch ty.Kind() {
    case reflect.Bool:

      return strconv.FormatBool(value.(bool))

    case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:

      return strconv.FormatInt(val.Int(), 10)

    case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:

      return strconv.FormatUint(val.Uint(), 10)

    case reflect.Float32, reflect.Float64:

      return strconv.FormatFloat(val.Float(), 'G', -1, 32)

    case reflect.Complex64, reflect.Complex128:

      return strconv.FormatComplex(val.Complex(), 'G', -1, 64)

    case reflect.String:

      return value.(string)
    }
  }

  return ""
}

func KStrRepr(value any) string {

  val := pp.KIndirectValueOf(value)

  if val.IsValid() {

    ty := val.Type()

    switch ty.Kind() {
    case reflect.String:

      // wrapping with quote
      return strconv.Quote(value.(string))
    }

    return KStrFmt(value)
  }

  // set zero value
  return pp.Noop[string]()
}
