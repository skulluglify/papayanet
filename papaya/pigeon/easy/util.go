package easy

import (
  "encoding/hex"
  "errors"
  "github.com/google/uuid"
  "reflect"
  "skfw/papaya/koala/mapping"
  "skfw/papaya/koala/pp"
  "strings"
)

func FieldGet(field any) (any, error) {

  val := pp.KIndirectValueOf(field)

  if val.IsValid() {

    return val.Interface(), nil
  }

  return nil, errors.New("unable to get value of field")
}

func FieldSet(field any, value any) error {

  var in reflect.Value

  in = pp.KIndirectValueOf(value)
  val := pp.KIndirectValueOf(field)

  if val.IsValid() && in.IsValid() {

    ty := val.Type()

    if ty.Kind() == in.Type().Kind() {

      if val.CanSet() {

        val.Set(in)
        return nil
      }
    }
  }

  return errors.New("unable to set value of field")
}

func MethodCall(method any, args ...any) []reflect.Value {

  var in []reflect.Value

  in = make([]reflect.Value, 0)
  val := pp.KIndirectValueOf(method)

  if val.IsValid() {

    ty := val.Type()

    switch ty.Kind() {
    case reflect.Func:

      for _, arg := range args {

        in = append(in, reflect.ValueOf(arg))
      }

      return val.Call(in)
    }
  }

  return nil
}

func TableName(model any) string {

  val := pp.KIndirectValueOf(model)

  if val.IsValid() {

    ty := val.Type()

    switch ty.Kind() {
    case reflect.Struct:
      method := val.MethodByName("TableName")
      out := MethodCall(method)
      if len(out) > 0 {

        // fallback default value
        return pp.Qstr(mapping.KValueToString(out[0]), "model")
      }
    }
  }

  return "model"
}

func StructGet(data any, key string) (any, error) {

  var err error
  var value any

  val := pp.KIndirectValueOf(data)

  if val.IsValid() {

    ty := val.Type()

    switch ty.Kind() {
    case reflect.Struct:

      method := val.FieldByName(key)

      if value, err = FieldGet(method); err != nil {

        return nil, err
      }

      return value, nil
    }
  }

  return nil, errors.New("invalid data structure")
}

func FindField(field reflect.Value, key string) reflect.Value {

  val := pp.KIndirectValueOf(field)

  if val.IsValid() {

    ty := val.Type()

    for i := 0; i < val.NumField(); i++ {

      f := pp.KIndirectValueOf(val.Field(i))
      ft := ty.Field(i)

      name := ft.Name

      if name == key {

        return f
      }

      switch f.Kind() {
      case reflect.Struct:

        field = FindField(f, key)

        if field.IsValid() {

          return field
        }

        break
      }
    }
  }

  return reflect.Value{}
}

func StructSet(data any, key string, value any) error {

  var err error

  val := pp.KIndirectValueOf(data)

  if val.IsValid() {

    ty := val.Type()

    switch ty.Kind() {
    case reflect.Struct:

      field := FindField(val, key)

      if field.IsValid() {

        if err = FieldSet(field, value); err != nil {

          return err
        }

        return nil
      }

      return errors.New("invalid data field")
    }
  }

  return errors.New("invalid data structure")
}

func EmptyIds(id string) bool {

  for _, c := range []byte(id) {

    if c != 48 { // 00 11 00 00

      return false
    }
  }

  return true
}

func EmptyId(id []byte) bool {

  return EmptyIds(hex.EncodeToString(id))
}

func EmptyIdx(id uuid.UUID) bool {

  var err error
  var idx []byte

  idx, err = id.MarshalBinary()

  if err != nil {

    return false
  }

  return EmptyId(idx)
}

func Id(id []byte) uuid.UUID {

  var err error
  var idx uuid.UUID

  idx, err = uuid.FromBytes(id)
  if err != nil {

    return uuid.UUID{}
  }

  return idx
}

func Idx(id uuid.UUID) string {

  var err error
  var data []byte

  data, err = id.MarshalBinary()
  if err != nil {

    return "00000000000000000000000000000000"
  }

  return hex.EncodeToString(data)
}

func Ids(id string) uuid.UUID {

  var err error
  var data []byte

  // remove all pad
  id = strings.ReplaceAll(id, "-", "")

  data, err = hex.DecodeString(id)
  if err != nil {

    return uuid.UUID{}
  }

  return Id(data)
}
