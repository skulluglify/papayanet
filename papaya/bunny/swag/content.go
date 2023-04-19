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

func SwagContentBoolean() m.KMapImpl {

	return &m.KMap{
		"type": "boolean",
	}
}

func SwagContentNumber() m.KMapImpl {

	return &m.KMap{
		"type": "number",
	}
}

func SwagContentString() m.KMapImpl {

	return &m.KMap{
		"type": "string",
	}
}

func SwagContentNullable() m.KMapImpl {

	return &m.KMap{
		"type": "null",
	}
}

func SwagContentArray(t m.KMapImpl) m.KMapImpl {

	return &m.KMap{
		"type":  "array",
		"items": t,
	}
}

func SwagContentObject(t m.KMapImpl) m.KMapImpl {

	return &m.KMap{
		"type":       "object",
		"properties": t,
	}
}

func SwagContentType(t string, v m.KMapImpl) m.KMapImpl {

	var cTy m.KMapImpl

	switch t {
	case "bool", "boolean":

		cTy = SwagContentBoolean()
		break

	case "int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64",
		"float", "float32", "float64",
		"complex", "complex64", "complex128",
		"integer", "decimal", "number", "byte": // byte as uint8

		cTy = SwagContentNumber()
		break

	case "str", "text", "string":

		cTy = SwagContentString()
		break

	case "array", "slice": // [] as slice

		if v != nil {

			cTy = SwagContentArray(v)
		}
		break

	case "map", "object":

		if v != nil {

			cTy = SwagContentObject(v)
		}
		break

	default:

		cTy = SwagContentNullable()
		break
	}

	return cTy
}

func SwagContentNormType(t string) string {

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

func SwagContentFormatter(mapping any) m.KMapImpl {

	var res m.KMapImpl
	val := pp.KIndirectValueOf(mapping)

	if val.IsValid() {

		ty := val.Type()

		switch ty.Kind() {
		case reflect.Array, reflect.Slice:

			if val.Len() > 0 {

				sample := val.Index(0).Interface()
				res = SwagContentArray(SwagContentFormatter(sample))

			} else {

				// sample type, normalize typing
				t := SwagContentNormType(ty.Elem().Name())

				// catch a typeof elem array or slice
				res = SwagContentArray(SwagContentType(t, nil))
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

					res = SwagContentObject(data)
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

			res = SwagContentObject(mm)
			break

		default:

			// type any in traditional typing, like bool, int, string
			res = SwagContentType(ty.Name(), nil)

			break
		}
	}

	return res
}
