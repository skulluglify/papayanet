package collection

import (
  "reflect"
  "skfw/papaya/koala/pp"
)

type Value struct {
  val reflect.Value
}

type ValueImpl interface {
  Init(value any)
  Compare(value any) CompareImpl[any]
  Bool() bool
  Int() int64
  Uint() uint64
  String() string
  Array() []any
  Map() map[string]any
  Any() any
}

func ValueNew(value any) ValueImpl {

  val := &Value{}
  val.Init(value)
  return val
}

func (v *Value) Init(value any) {

  v.val = pp.KIndirectValueOf(value)
}

func (v *Value) Compare(value any) CompareImpl[any] {

  val := pp.KIndirectValueOf(value)
  return CompareNew[any](v.val, val)
}

func (v *Value) Bool() bool {

  val := v.val

  if val.IsValid() {

    ty := val.Type()

    switch ty.Kind() {
    case reflect.Bool:

      return val.Bool()
    }
  }

  return false
}

func (v *Value) Int() int64 {

  val := v.val

  if val.IsValid() {

    ty := val.Type()

    switch ty.Kind() {
    case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:

      return val.Int()
    }
  }

  return 0
}

func (v *Value) Uint() uint64 {

  val := v.val

  if val.IsValid() {

    ty := val.Type()

    switch ty.Kind() {
    case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:

      return val.Uint()
    }
  }

  return 0
}

func (v *Value) String() string {

  val := v.val

  if val.IsValid() {

    ty := val.Type()

    switch ty.Kind() {
    case reflect.String:

      return val.String()
    }
  }

  return ""
}

func (v *Value) Array() []any {

  val := v.val

  if val.IsValid() {

    ty := val.Type()

    switch ty.Kind() {
    case reflect.Array, reflect.Slice:

      var res []any

      for i := 0; i < val.Len(); i++ {

        res = append(res, val.Index(i).Interface())
      }

      return res
    }
  }

  return nil
}

func (v *Value) Map() map[string]any {

  val := v.val

  if val.IsValid() {

    ty := val.Type()

    switch ty.Kind() {
    case reflect.Map:

      if ty.Key().Kind() == reflect.String {

        res := map[string]any{}
        iter := val.MapRange()

        for iter.Next() {

          key, value := iter.Key(), iter.Value()
          res[key.String()] = value.Interface()
        }

        return res
      }
    }
  }

  return nil
}

func (v *Value) Any() any {

  val := v.val

  if val.IsValid() {

    return val.Interface()
  }

  return nil
}
