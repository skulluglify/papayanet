package environ

import (
  "os"
  "reflect"
  "strconv"
)

type KEnvLoader[T any] struct{}
type KEnvLoaderImpl[T any] interface {
  Load(value T) bool
}

func KEnvLoaderNew[T any]() KEnvLoaderImpl[T] {

  envLoader := &KEnvLoader[T]{}
  return envLoader
}

func (e *KEnvLoader[T]) Load(value T) bool {

  val := reflect.Indirect(reflect.ValueOf(value))

  if val.IsValid() {

    ty := val.Type()

    switch ty.Kind() {
    case reflect.Struct:

      var i, n int
      var vf reflect.Value
      var vfTy reflect.Type
      var tf reflect.StructField
      var tags reflect.StructTag
      //var sName, sValue, sTag, sEnv string
      var sTag, sEnv string

      n = ty.NumField()

      for i = 0; i < n; i++ {

        tf, vf = ty.Field(i), val.Field(i)
        tags = tf.Tag

        if tf.IsExported() {

          //sName = tf.Name
          //sValue = vf.String()
          sTag = tags.Get("env")

          if sTag != "" {

            vfTy = tf.Type
            sEnv = os.Getenv(sTag)

            if sEnv != "" {

              switch vfTy.Kind() {
              case reflect.Int:

                // skipping cvt error, set zero value
                res, _ := strconv.ParseInt(sEnv, 10, strconv.IntSize)
                vf.Set(reflect.ValueOf(int(res)))
                continue

              case reflect.Uint:

                // skipping cvt error, set zero value
                res, _ := strconv.ParseUint(sEnv, 10, strconv.IntSize)
                vf.Set(reflect.ValueOf(uint(res)))
                continue

              case reflect.Bool:

                // skipping cvt error, set zero value
                res, _ := strconv.ParseBool(sEnv)
                vf.Set(reflect.ValueOf(res))
                continue

              case reflect.String:

                vf.Set(reflect.ValueOf(sEnv))
                continue
              }
            }
          }
        }

        // missing environ
        //return false
      }

      // all set up, don't care about missing environ
      return true
    }
  }

  // not valid
  return false
}
