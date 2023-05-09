package pp

import (
	"reflect"
)

// Method for determine Coalescing Operator

func Q[T any](values ...T) T {

	// catch valid value
	for _, value := range values {

		val := KIndirectValueOf(value)

		// `nil` has been handled it
		if val.IsValid() {

			ty := val.Type()
			switch ty.Kind() {
			case reflect.Bool:

				if !val.Bool() {

					continue
				}

				break

			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:

				if val.Int() == 0 {

					continue
				}

				break

			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:

				if val.Uint() == 0 {

					continue
				}

				break

			case reflect.Float32, reflect.Float64:

				if val.Float() == 0.0 {

					continue
				}

				break

			case reflect.Complex64, reflect.Complex128:

				if val.Complex() == 0 {

					continue
				}

				break

			case reflect.String:

				if val.String() == "" {

					continue
				}

				break
			}

			return value
		}
	}

	// zero value, as a default value
	return Noop[T]()
}

// auto type defined by name

var Qany = Q[any]

// make it fast implementation

func Qstr(values ...string) string {

	for _, value := range values {

		if value != "" {

			return value
		}
	}

	return Noop[string]()
}

func Qbool(values ...bool) bool {

	for _, value := range values {

		if value {

			return value
		}
	}

	return Noop[bool]()
}

func Qbyte(values ...byte) byte {

	for _, value := range values {

		if value != 0 {

			return value
		}
	}

	return Noop[byte]()
}

func Qint(values ...int) int {

	for _, value := range values {

		if value != 0 {

			return value
		}
	}

	return Noop[int]()
}

func Quint(values ...uint) uint {

	for _, value := range values {

		if value != 0 {

			return value
		}
	}

	return Noop[uint]()
}

// ---

// Method for validity Flag is true

func Qflag(flags int, flag int) bool {

	return flags|flag != 0 // or flag == 1
}
