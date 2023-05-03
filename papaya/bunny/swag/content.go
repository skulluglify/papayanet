package swag

import (
  "reflect"
  m "skfw/papaya/koala/mapping"
  "skfw/papaya/koala/pp"
)

func SwagContentSchema(mimeType string, data any, description string) m.KMapImpl {

  content := &m.KMap{}

  switch mimeType {
  case "application/json", "application/xml", "multipart/form-data":

    content.Put(mimeType, &m.KMap{
      "schema": SwagContentFormatter(data),
    })

    break

  default:

    // default media type is binary
    content.Put(mimeType, &m.KMap{
      "schema": &m.KMap{
        "type":   "string",
        "format": "binary",
      },
    })

    break

  }

  return &m.KMap{
    "description": description,
    "content":     content,
  }
}

func SwagContentSchemes(body m.KMapImpl) []m.KMapImpl {

  res := make([]m.KMapImpl, 0)
  if mm := m.KMapCast(body); mm != nil {

    for _, enum := range mm.Enums() {

      k, v := enum.Tuple()

      if mimeTy := m.KValueToString(k); mimeTy != "" {

        if vM := m.KMapCast(v); vM != nil {

          schema := vM.Get("schema")
          description := pp.QStr(m.KValueToString(vM.Get("description")), "Ok")
          res = append(res, SwagContentSchema(mimeTy, schema, description))
        }
      }

    }
  }

  return res
}

func SwagContentFormatter(mapping any) m.KMapImpl {

  var res m.KMapImpl
  val := pp.KIndirectValueOf(mapping)

  if val.IsValid() {

    ty := val.Type()

    switch ty.Kind() {
    case reflect.Array, reflect.Slice:

      if val.Len() > 0 {

        sample := val.Index(0).Interface()
        res = SwagUniversalArray(SwagContentFormatter(sample))

      } else {

        // sample type, normalize typing
        t := SwagUniversalNormType(ty.Elem().Name())

        // catch a typeof elem array or slice
        res = SwagUniversalArray(SwagUniversalType(t, nil))
      }

      break

    case reflect.Map:

      if ty == reflect.TypeOf(m.KMap{}) {

        sample := val.Interface()
        if mm := m.KMapCast(sample); mm != nil {

          data := &m.KMap{}

          for _, enum := range mm.Enums() {

            k, v := enum.Tuple()

            data.Put(k, SwagContentFormatter(v))
          }

          res = SwagUniversalObject(data)
        }
      }

      break

    case reflect.Struct:

      var i, n int
      var vf reflect.Value
      var vt reflect.StructField
      var name, tag string
      var value any

      n = val.NumField()

      // convert struct as object mapping
      mm := &m.KMap{}

      for i = 0; i < n; i++ {

        vf, vt = val.Field(i), ty.Field(i)

        if vt.IsExported() {

          if vf.IsValid() {

            name = vt.Name
            tag = vt.Tag.Get("json")
            value = vf.Interface()
            if tag != "" {

              name = tag
            }

            // put magic
            mm.Put(name, SwagContentFormatter(value))
          }
        }
      }

      res = SwagUniversalObject(mm)
      break

    case reflect.String:

      // maybe text type string

      t := val.String()

      var retype bool
      retype = false

      if t != "" {

        if m.Keys(Types).Contain(t) {

          res = SwagUniversalType(t, nil)
          retype = true
        }
      }

      // fallback use type string as default
      if !retype {

        res = SwagUniversalType(ty.Name(), nil)
      }

      break

    default:

      // type any in traditional typing, like bool, int, string
      res = SwagUniversalType(ty.Name(), nil)

      break
    }
  }

  return res
}
