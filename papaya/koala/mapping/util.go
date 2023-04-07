package mapping

import (
  "PapayaNet/papaya/koala"
  "PapayaNet/papaya/koala/pp"
  "reflect"
  "strconv"
  "strings"
)

func KMapEnums(mapping any) koala.KEnums[string, any] {

  enums := koala.KEnumsNew[string, any](0)

  val := pp.KIndirectValueOf(mapping)

  if val.IsValid() {

    ty := val.Type()

    switch ty.Kind() {
    case reflect.Array, reflect.Slice:

      for i := 0; i < val.Len(); i++ {

        k, v := strconv.Itoa(i), val.Index(i).Interface()
        enum := koala.KEnumNew[string, any](k, v)
        enums = append(enums, enum)
      }

      break

    case reflect.Map:

      if ty.Key().Kind() == reflect.String {

        mapIter := val.MapRange()

        for mapIter.Next() {

          k, v := mapIter.Key().String(), mapIter.Value().Interface()
          enum := koala.KEnumNew[string, any](k, v)
          enums = append(enums, enum)
        }
      }

      break
    case reflect.Struct:

      // TODO: not implemented yet
      var i, n int
      var vf reflect.Value
      var tf reflect.StructField

      n = val.NumField()

      for i = 0; i < n; i++ {

        tf, vf = ty.Field(i), val.Field(i)

        if tf.IsExported() {

          if vf.IsValid() {

            k, v := tf.Name, vf.Interface()
            enum := koala.KEnumNew[string, any](k, v)
            enums = append(enums, enum)
          }
        }
      }

      break
    }
  }

  return enums
}

func KMapTreeEnums(mapping any) koala.KEnums[string, any] {

  enums := koala.KEnumsNew[string, any](0)

  for _, first := range KMapEnums(mapping) {

    k, v := first.Tuple()

    enum := koala.KEnumNew[string, any](k, v)
    enums = append(enums, enum)

    for _, second := range KMapTreeEnums(v) {

      kk, vv := second.Tuple()
      enum := koala.KEnumNew[string, any](k+"."+kk, vv)
      enums = append(enums, enum)
    }
  }

  return enums
}

func KMapKeys(mapping any) []string {

  tokens := make([]string, 0)

  val := pp.KIndirectValueOf(mapping)

  if val.IsValid() {

    ty := val.Type()

    switch ty.Kind() {
    case reflect.Array, reflect.Slice:

      for i := 0; i < val.Len(); i++ {

        k := strconv.Itoa(i)
        tokens = append(tokens, k)
      }

      break

    case reflect.Map:

      if ty.Key().Kind() == reflect.String {

        mapIter := val.MapRange()

        for mapIter.Next() {

          k := mapIter.Key().String()
          tokens = append(tokens, k)
        }
      }

      break
    case reflect.Struct:

      // TODO: not implemented yet
      var i, n int
      var vf reflect.Value
      var tf reflect.StructField

      n = val.NumField()

      for i = 0; i < n; i++ {

        tf, vf = ty.Field(i), val.Field(i)

        if tf.IsExported() {

          if vf.IsValid() {

            k := tf.Name
            tokens = append(tokens, k)
          }
        }
      }

      break
    }
  }

  return tokens
}

func KMapTreeKeys(mapping any) []string {

  tokens := make([]string, 0)

  for _, enum := range KMapEnums(mapping) {

    k, v := enum.Tuple()

    tokens = append(tokens, k)

    for _, kk := range KMapTreeKeys(v) {

      tokens = append(tokens, k+"."+kk)
    }
  }

  return tokens
}

func KMapValues(mapping any) []any {

  values := make([]any, 0)

  val := pp.KIndirectValueOf(mapping)

  if val.IsValid() {

    ty := val.Type()

    switch ty.Kind() {
    case reflect.Array, reflect.Slice:

      for i := 0; i < val.Len(); i++ {

        v := val.Index(i)
        values = append(values, v)
      }

      break

    case reflect.Map:

      if ty.Key().Kind() == reflect.String {

        mapIter := val.MapRange()

        for mapIter.Next() {

          v := mapIter.Value().Interface()
          values = append(values, v)
        }
      }

      break
    case reflect.Struct:

      // TODO: not implemented yet
      var i, n int
      var vf reflect.Value
      var tf reflect.StructField

      n = val.NumField()

      for i = 0; i < n; i++ {

        tf, vf = ty.Field(i), val.Field(i)

        if tf.IsExported() {

          if vf.IsValid() {

            v := vf.Interface()
            values = append(values, v)
          }
        }
      }

      break
    }
  }

  return values
}

// TODO: Branch, and Put, not stable yet
// TODO: Branch, and Put, not enough implemented yet, like Array, Slice, Struct, and other

func KMapBranch(name string, mapping any) any {

  val := reflect.ValueOf(mapping)
  tokens := strings.Split(name, ".")

  if val.IsValid() {

    var currKey, prevKey string

    for i, k := range tokens {

      // new `val` is invalid
      if !val.IsValid() {

        break
      }

      if i == 0 {

        currKey = k

      } else {

        currKey += "." + k
      }

      val = pp.KIndirectValueOf(val)
      ty := val.Type()

      // search

      m := KMapGetValue(k, val.Interface())

      // replace

      if m != nil {

        val = reflect.ValueOf(m)
        prevKey = currKey
        continue
      }

      // create & set on previous existing
      if prevKey != "" {

        mm := &KMap{}
        if KMapSetValue(prevKey, mm, val.Interface()) {

          val = reflect.ValueOf(mm)
          prevKey = currKey
          continue
        }

      }

      // hacked, put new assign in a map object
      switch ty.Kind() {
      case reflect.Map:
        if ty.Key().Kind() == reflect.String {

          if ty == reflect.TypeOf(KMap{}) {

            mm := &KMap{}
            t := val.Interface()
            o := map[string]any(t.(KMap))
            o[k] = mm

            val = reflect.ValueOf(mm)
            prevKey = currKey
            continue
          }
        }

        break

      default: // no assignable

        return nil
      }

      prevKey = currKey
    }

    // passing, that right or not, on mystery
    return val
  }

  return nil
}

func KMapPut(name string, data any, mapping any) bool {

  val := pp.KIndirectValueOf(mapping)

  if val.IsValid() {

    ty := val.Type()

    if !KMapSetValue(name, data, mapping) {

      var prefix, suffix string
      tokens := strings.Split(name, ".")
      n := len(tokens)

      if n > 0 {

        if n > 1 {

          prefix = strings.Join(tokens[:n-1], ".")
          suffix = tokens[n-1]

        } else {

          prefix = tokens[0]
        }
      }

      if prefix != "" {

        if suffix != "" {

          // check and build new branch
          if m := KMapBranch(prefix, mapping); m != nil {

            return KMapPut(suffix, data, m)
          }

          return false
        }

        // hacked, put new assign in a map object
        switch ty.Kind() {
        case reflect.Map:

          if ty.Key().Kind() == reflect.String {

            if ty == reflect.TypeOf(KMap{}) {

              mm := map[string]any(val.Interface().(KMap))
              mm[prefix] = data
              return true
            }
          }
        }
      }

      return false
    }
  }

  return true
}
