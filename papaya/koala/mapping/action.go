package mapping

import (
	"PapayaNet/papaya/koala/pp"
	"reflect"
	"strconv"
	"strings"
)

func KMapGetValue(name string, mapping any) any {

	tokens := strings.Split(name, ".")

	var value any
	value = mapping

	i := 0
	n := len(tokens)

	if n == 0 {

		return nil
	}

	for {

		if value == nil {

			return nil
		}

		if n <= i {

			break
		}

		token := tokens[i]

		val := pp.KIndirectValueOf(value)

		if val.IsValid() {

			ty := val.Type()

			switch ty.Kind() {

			case reflect.Array, reflect.Slice:

				if index, err := strconv.Atoi(token); err == nil {

					if index < val.Len() {

						v := val.Index(index)
						value = v.Interface()
						break
					}
				}

				value = nil

				break

			case reflect.Map:

				if ty.Key().Kind() == reflect.String {

					m := false

					mapIter := val.MapRange()

					for mapIter.Next() {

						key := mapIter.Key().String()

						if token == key {

							//value := reflect.Indirect(mapIter.Value())
							v := mapIter.Value()
							value = v.Interface()
							m = true
							break
						}
					}

					if !m {

						value = nil
					}
				}

				break

			case reflect.Struct:

				var m bool
				var j, k int
				var vf reflect.Value
				var tf reflect.StructField

				m, k = false, val.NumField()

				for j = 0; j < k; j++ {

					tf, vf = ty.Field(j), val.Field(j)

					if tf.IsExported() {

						if vf.IsValid() {

							if token == tf.Name {

								value = vf.Interface()
								m = true
								break
							}
						}
					}
				}

				if !m {

					value = nil
				}

				break
			}
		}

		i++
	}

	return value
}

func KMapSetValue(name string, data any, mapping any) bool {

	tokens := strings.Split(name, ".")

	value := reflect.ValueOf(mapping)

	i := 0
	n := len(tokens)

	if n == 0 {

		return false
	}

	// set value as a pointer ? im don't thing so, because can set it without catch pointer
	// `array`, `slice`, or `map` that combine ptr on inside memory
	// that get previous and set value

	for {

		if !value.IsValid() {

			return false
		}

		if n <= i {

			break
		}

		token := tokens[i]

		// that have problem, type hide on inside interface and ptr
		// how to solve, is catch all elems on inside interface and set back in var `value`
		// a reset type in var `ty`
		value = pp.KIndirectValueOf(value)

		// update type of temporary

		if value.IsValid() {

			ty := value.Type()

			// lookup data on `map`
			switch value.Kind() {

			case reflect.Array, reflect.Slice:

				if index, err := strconv.Atoi(token); err == nil {

					if index < value.Len() {

						// get previous, set value
						if i+1 == n {

							// get previous, to set value on inside `array` or `slice`
							value.Index(index).Set(reflect.ValueOf(data))
							return true
						}

						//v := value.Index(index)
						//value = v
						value = value.Index(index)
						break
					}

					// index out of bound
					return false
				}

				// NaN
				return false

			case reflect.Map:
				if ty.Key().Kind() == reflect.String {

					m := false

					mapIter := value.MapRange()

					for mapIter.Next() {

						key := mapIter.Key()

						if token == key.String() {

							// get previous, set value
							if i+1 == n {

								// get previous, to set value on inside `map`
								value.SetMapIndex(key, reflect.ValueOf(data))
								return true
							}

							//v := mapIter.Value()
							//value = v
							value = mapIter.Value()

							m = true
							break
						}
					}

					if !m {

						return false
					}

					break
				}

				// bad key
				return false

			case reflect.Struct:

				var m bool
				var j, k int
				var vf reflect.Value
				var tf reflect.StructField

				m, k = false, value.NumField()

				for j = 0; j < k; j++ {

					tf, vf = ty.Field(j), value.Field(j)

					if tf.IsExported() {

						if token == tf.Name {

							if i+1 == n {

								vf.Set(reflect.ValueOf(data))
								return true
							}

							//v := vf
							//value = v
							value = vf

							m = true
							break
						}
					}
				}

				if !m {

					return false
				}

				break
			}
		}

		i++
	}

	// other than `array`, `slice`, or `map`
	// can't set value as ptr
	return false
}

func KMapDelValue(name string, mapping any) bool {

	tokens := strings.Split(name, ".")

	var value, prev reflect.Value

	value = reflect.ValueOf(mapping)
	prev = reflect.Value{}

	var prevToken string

	var i, n int
	i, n = 0, len(tokens)

	if n == 0 {

		return false
	}

	// set value as a pointer ? im don't thing so, because can set it without catch pointer
	// `array`, `slice`, or `map` that combine ptr on inside memory
	// that get previous and set value

	for {

		if !value.IsValid() {

			return false
		}

		if n <= i {

			break
		}

		token := tokens[i]

		// that have problem, type hide on inside interface and ptr
		// how to solve, is catch all elems on inside interface and set back in var `value`
		// a reset type in var `ty`
		value = pp.KIndirectValueOf(value)

		if value.IsValid() {

			ty := value.Type()

			// lookup data on `map`
			switch value.Kind() {

			case reflect.Array, reflect.Slice:

				if index, err := strconv.Atoi(token); err == nil {

					if index < value.Len() {

						// get previous, set value
						if i+1 == n {

							s := value.Len()

							// get previous, to set value on inside `array` or `slice`
							L := value.Slice(0, index)
							R := value.Slice(index+1, s)

							//k := L.Len() + R.Len()

							// ERROR: problem
							//data := reflect.MakeSlice(ty, k, k)
							//
							//// merging
							//// TODO: out of time, make it fast
							//for j := 0; j < L.Len(); j++ {
							//
							//	data.Index(j).Set(L.Index(j))
							//}
							//
							//for j := 0; j < R.Len(); j++ {
							//
							//	data.Index(j + L.Len()).Set(R.Index(j))
							//}

							// make it fast
							data := reflect.AppendSlice(L, R)

							// end

							// save on previous value as ref
							if reflect.ValueOf(prevToken).IsValid() {

								switch prev.Kind() {

								case reflect.Array, reflect.Slice:

									if j, e := strconv.Atoi(prevToken); e == nil {

										prev.Index(j).Set(data)
										return true
									}

									return false

								case reflect.Map:

									prev.SetMapIndex(reflect.ValueOf(prevToken), data)
									break

								case reflect.Struct:

									prev.FieldByName(prevToken).Set(data)
									break
								}

								return true
							}

							// try save on current elem
							// panic: reflect: reflect.Value.Set using unaddressable value
							//value.Set(data)
							//if reflect.DeepEqual(data, value.Interface()) {
							//
							//  return true
							//}

							return false
						}

						//v := value.Index(index)
						//prev = value
						//value = v
						prev = value
						value = value.Index(index)

						break
					}

					// index out of bound
					return false
				}

				// NaN
				return false

			case reflect.Map:

				if ty.Key().Kind() == reflect.String {

					m := false

					mapIter := value.MapRange()

					for mapIter.Next() {

						key := mapIter.Key()

						if token == key.String() {

							// get previous, set value
							if i+1 == n {

								// get previous, to delete value on inside `map`
								value.SetMapIndex(key, reflect.Value{}) // that don't know how it works
								return true
							}

							//v := mapIter.Value()
							//prev = value
							//value = v
							prev = value
							value = mapIter.Value()

							m = true
							break
						}
					}

					if !m {

						return false
					}

					break
				}

				// bad key
				return false

			case reflect.Struct:

				var m bool
				var j, k int
				var vf reflect.Value
				var tf reflect.StructField

				m, k = false, value.NumField()

				for j = 0; j < k; j++ {

					tf, vf = ty.Field(j), value.Field(j)

					if tf.IsExported() {

						// delete action, require value is valid
						if vf.IsValid() {

							if token == tf.Name {

								if i+1 == n {

									// make it zero value
									vf.Set(reflect.Zero(tf.Type))
									return true
								}

								//v := vf
								//prev = value
								//value = v
								prev = value
								value = vf

								m = true
								break
							}
						}
					}
				}

				if !m {

					return false
				}

				break
			}
		}

		prevToken = token
		i++
	}

	// other than `array`, `slice`, or `map`
	// can't set value as ptr
	return false
}
