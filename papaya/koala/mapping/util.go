package mapping

import (
	"reflect"
	"skfw/papaya/koala/collection"
	"skfw/papaya/koala/pp"
	"strconv"
)

func KMapEnums(mapping any) []collection.KEnumImpl[string, any] {

	enums := make([]collection.KEnumImpl[string, any], 0)

	val := pp.KIndirectValueOf(mapping)

	if val.IsValid() {

		ty := val.Type()

		switch ty.Kind() {
		case reflect.Array, reflect.Slice:

			for i := 0; i < val.Len(); i++ {

				k, v := strconv.Itoa(i), val.Index(i).Interface()
				enum := collection.KEnumNew[string, any](k, v)
				enums = append(enums, enum)
			}

			break

		case reflect.Map:

			if ty.Key().Kind() == reflect.String {

				mapIter := val.MapRange()

				for mapIter.Next() {

					k, v := mapIter.Key().String(), mapIter.Value().Interface()
					enum := collection.KEnumNew[string, any](k, v)
					enums = append(enums, enum)
				}
			}

			break
		case reflect.Struct:

			// TODO: not implemented yet
			var i, n int
			var vf reflect.Value
			var tags reflect.StructTag
			var tf reflect.StructField

			n = val.NumField()

			for i = 0; i < n; i++ {

				tf, vf = ty.Field(i), val.Field(i)
				tags = tf.Tag

				if tf.IsExported() {

					if vf.IsValid() {

						k, v := tf.Name, vf.Interface()
						t := pp.Qstr(tags.Get("tag"), tags.Get("json"))
						enum := collection.KEnumNew[string, any](pp.Qstr(t, k), v)
						enums = append(enums, enum)
					}
				}
			}

			break
		}
	}

	return enums
}

func KMapTreeEnums(mapping any) []collection.KEnumImpl[string, any] {

	enums := make([]collection.KEnumImpl[string, any], 0)

	for _, first := range KMapEnums(mapping) {

		k, v := first.Tuple()

		enum := collection.KEnumNew[string, any](k, v)
		enums = append(enums, enum)

		for _, second := range KMapTreeEnums(v) {

			kk, vv := second.Tuple()
			enum := collection.KEnumNew[string, any](k+"."+kk, vv)
			enums = append(enums, enum)
		}
	}

	return enums
}

func KMapKeys(mapping any) Keys {

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
			var tags reflect.StructTag
			var tf reflect.StructField

			n = val.NumField()

			for i = 0; i < n; i++ {

				tf, vf = ty.Field(i), val.Field(i)
				tags = tf.Tag

				if tf.IsExported() {

					if vf.IsValid() {

						k := tf.Name
						t := pp.Qstr(tags.Get("tag"), tags.Get("json"))
						tokens = append(tokens, pp.Qstr(t, k))
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

// try casting into KMap

func KMapCast(mapping any) KMapImpl {

	// catch from interface or pointer
	val := pp.KIndirectValueOf(mapping)

	// valid value is not null
	if val.IsValid() {

		if mm, ok := val.Interface().(KMap); ok {

			return &mm
		}
	}

	return nil
}
