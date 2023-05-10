package gen

import (
  "reflect"
  "skfw/papaya/koala/pp"
  "strconv"
  "strings"
)

type KMapIteration struct {
  KIterationImpl[string, any]
}

type KMapIterationNextHandler KIterationNextHandler[string, any]

type KMapIterationImpl interface {
  KIterationImpl[string, any]
}

func KMapStopIteration() KMapIterationImpl {

  return &KMapIteration{
    &KIteration[string, any]{
      stopIter: true,
    },
  }
}

func KMapIterable(mapping any) KMapIterationImpl {

  val := pp.KIndirectValueOf(mapping)

  if val.IsValid() {

    ty := val.Type()

    switch ty.Kind() {
    case reflect.Array, reflect.Slice:

      n := val.Len()
      k := 0

      if n > 0 {

        return &KMapIteration{
          &KIteration[string, any]{
            NextHandler: func(v KIterationImpl[string, any]) error {

              if k < n {

                key := strconv.Itoa(k)
                value := val.Index(k).Interface()

                v.SetValues(key, value)
              }

              k += 1

              return nil
            },
          },
        }
      }

      break
    case reflect.Map:

      if ty.Key().Kind() == reflect.String {

        mapRange := val.MapRange()

        return &KMapIteration{
          &KIteration[string, any]{
            NextHandler: func(v KIterationImpl[string, any]) error {

              if hasNext := mapRange.Next(); hasNext {

                key := mapRange.Key().String()
                value := mapRange.Value().Interface()

                v.SetValues(key, value)
              }

              return nil
            },
          },
        }
      }

      break

    case reflect.Struct:

      var k, n int
      var vf reflect.Value
      var tf reflect.StructField
      var tags reflect.StructTag
      var sName, sTag string
      //var sName string
      var sValue any

      k, n = 0, ty.NumField()

      if n > 0 {

        return &KMapIteration{
          &KIteration[string, any]{
            NextHandler: func(v KIterationImpl[string, any]) error {

              // wait for search field is exported
              // example
              // k 0 is exported
              // k 1 is not exported
              // k 2 is not exported
              // k 3 is exported
              // k 4 closed
              for {

                if n <= k {

                  break
                }

                tf, vf = ty.Field(k), val.Field(k)
                tags = tf.Tag

                if tf.IsExported() {

                  if vf.IsValid() {

                    sName = tf.Name
                    sValue = vf.Interface()
                    sTag = pp.Qstr(tags.Get("tag"), tags.Get("json"))

                    v.SetValues(pp.Qstr(sTag, sName), sValue)
                    k += 1
                    break
                  }
                }

                k += 1
              }

              return nil
            },
          },
        }
      }
    }
  }

  return KMapStopIteration()
}

type KMapPageIteration struct {
  KPageIterationImpl[string, any]
}

type KMapPageIterationImpl interface {
  KPageIterationImpl[string, any]
}

func KMapTreeIterable(mapping any) KMapIterationImpl {

  //val := koala.KIndirectValueOf(mapping)
  //ty := val.Type()

  mapPageIteration := &KMapPageIteration{
    &KPageIteration[string, any]{},
  }
  mapPageIteration.Init()

  iter := KMapIterable(mapping)
  mapPageIteration.Add(iter)

  //k := 0

  return &KMapIteration{
    &KIteration[string, any]{
      NextHandler: func(v KIterationImpl[string, any]) error {

        next := mapPageIteration.Wait()

        //// cut off
        //if k > 6000 {
        //  return v.Stop()
        //}
        //k += 1

        if next.HasNext() {

          enum := next.Enum()
          value := enum.Value()
          keys := mapPageIteration.Keys()

          if KMapHunt(value) {

            iter := KMapIterable(value)
            mapPageIteration.Add(iter)
          }

          v.SetValues(strings.Join(keys, "."), value)
        }

        return nil
      },
    },
  }
}
