package pp

import "reflect"

// hack valueOf from `reflect.Value` to get indirect as value
// read as interface ptr, ptr, and passing
// better than `reflect.Indirect`

func KIndirectValueOf(data any) reflect.Value {

  var val reflect.Value

  //if reflect.DeepEqual(
  //  reflect.TypeOf(data),
  //  reflect.TypeOf(reflect.Value{})) {

  if reflect.TypeOf(data) == reflect.TypeOf(reflect.Value{}) {

    val = data.(reflect.Value)

  } else {

    val = reflect.ValueOf(data)
  }

  // safety
  if val.IsValid() {

    ty := val.Type()

    switch ty.Kind() {

    case reflect.Interface, reflect.Pointer: // Catch Any Type

      // unsafe, make infinity loop
      // recursion if contain interface, or ptr again
      return KIndirectValueOf(val.Elem())
    }
  }

  return val
}
