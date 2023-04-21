package swag

import (
  m "PapayaNet/papaya/koala/mapping"
  "PapayaNet/papaya/koala/pp"
  "reflect"
  "strings"
)

// convert expect into openapi format

// boolean
// number
// string
// array
// object
// null

func SwagUniversalBoolean() m.KMapImpl {

  return &m.KMap{
    "type": "boolean",
  }
}

func SwagUniversalNumber() m.KMapImpl {

  return &m.KMap{
    "type": "number",
  }
}

func SwagUniversalString() m.KMapImpl {

  return &m.KMap{
    "type": "string",
  }
}

func SwagUniversalNull() m.KMapImpl {

  return &m.KMap{
    "type": "null",
  }
}

func SwagUniversalArray(t m.KMapImpl) m.KMapImpl {

  return &m.KMap{
    "type":  "array",
    "items": t,
  }
}

func SwagUniversalObject(t m.KMapImpl) m.KMapImpl {

  return &m.KMap{
    "type":       "object",
    "properties": t,
  }
}

func SwagUniversalReType(v any) string {

  val := pp.KIndirectValueOf(v)

  if val.IsValid() {

    ty := val.Type()

    switch ty.Kind() {

    case reflect.Bool:
      return "bool"

    case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
      reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
      reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:

      return "number"

    case reflect.Array, reflect.Slice:

      return "array"

    case reflect.Map, reflect.Struct:

      return "object"

    case reflect.String:
      
      return pp.QStr(m.KValueToString(SwagUniversalType(val.String(), nil).Get("type")), "string")
    }
  }

  return "null"
}

func SwagUniversalType(t string, v m.KMapImpl) m.KMapImpl {

  var cTy m.KMapImpl

  switch t {
  case "bool", "boolean":

    cTy = SwagUniversalBoolean()
    break

  case "int", "int8", "int16", "int32", "int64",
    "uint", "uint8", "uint16", "uint32", "uint64",
    "float", "float32", "float64",
    "complex", "complex64", "complex128",
    "integer", "decimal", "number", "byte": // byte as uint8

    cTy = SwagUniversalNumber()
    break

  case "str", "text", "string":

    cTy = SwagUniversalString()
    break

  case "array", "slice": // [] as slice

    if v != nil {

      cTy = SwagUniversalArray(v)
    }
    break

  case "map", "object":

    if v != nil {

      cTy = SwagUniversalObject(v)
    }
    break

  default:

    //case "nil", "null":
    cTy = SwagUniversalNull()
    break
  }

  return cTy
}

// array, slice cases
// can be map, array, or slice again

func SwagUniversalNormType(t string) string {

  // map.+? is map
  if strings.HasSuffix(t, "map") {

    return "map"
  }

  // [].+? is slice
  if strings.HasSuffix(t, "[]") {

    return "slice"
  }

  // [.+? is array
  if strings.HasSuffix(t, "[") {

    return "array"
  }

  return t // as null
}
