package mapping

import (
	"PapayaNet/papaya/koala/pp"
	"reflect"
	"strconv"
	"strings"
)

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

			// some good point for a set array like slicing, or append
			// create & set on previous existing
			//if prevKey != "" {
			//
			//	mm := &KMap{}
			//	if KMapSetValue(prevKey, mm, val.Interface()) {
			//
			//		val = reflect.ValueOf(mm)
			//		prevKey = currKey
			//		continue
			//	}
			//
			//}

			// var `k` may KMap or List, whatever, not implemented yet

			// hacked, put new assign in a map object
			switch ty.Kind() {
			case reflect.Array, reflect.Slice:

				if data := KMapGetValue(prevKey, mapping); data != nil {

					v := pp.KIndirectValueOf(data)

					if v.IsValid() {

						switch v.Type().Kind() {
						case reflect.Array, reflect.Slice:

							if index, err := strconv.Atoi(k); err == nil {

								mm := &KMap{}
								if index < v.Len() {

									v.Index(index).Set(reflect.ValueOf(mm))
									break
								}

								// try
								if o, ok := data.([]any); ok {

									o = append(o, mm)
									KMapSetValue(prevKey, mapping, o)
								}
							}

							break
						}
					}

				}

				panic("not stable yet")

			case reflect.Map:
				if ty.Key().Kind() == reflect.String {

					if ty == reflect.TypeOf(KMap{}) {

						mm := &KMap{}
						v := val.Interface()
						o := map[string]any(v.(KMap))
						o[k] = mm

						val = reflect.ValueOf(mm)
						prevKey = currKey
						continue
					}
				}

				break

				// skipping `struct`

			default: // no assignable

				return nil
			}

			prevKey = currKey
		}

		// void keep
		func(_ string) {}(prevKey)

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

							v := val.Interface()
							mm := map[string]any(v.(KMap))
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
