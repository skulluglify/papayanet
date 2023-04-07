package koala

import (
	"PapayaNet/papaya/panda"
	"reflect"
	"strconv"
	"strings"
)

type KMap map[string]any

func KMapGetValue(name string, mapping any) any {

	tokens := strings.Split(name, ".")

	var temp any
	temp = mapping

	i := 0
	n := len(tokens)

	if n == 0 {

		return nil
	}

	for {

		if temp == nil {

			return nil
		}

		if n <= i {
			break
		}

		token := tokens[i]

		v := reflect.Indirect(reflect.ValueOf(temp))
		t := v.Type()

		switch t.Kind() {

		case reflect.Array, reflect.Slice:

			if index, err := strconv.Atoi(token); err == nil {

				if index < v.Len() {

					//value := reflect.Indirect(v.Index(index))
					value := v.Index(index)
					temp = value.Interface()
					break
				}
			}

			// not found
			temp = nil

			break

		case reflect.Map:
			if t.Key().Kind() == reflect.String {

				m := false

				mapIter := v.MapRange()

				for mapIter.Next() {

					key := mapIter.Key().String()

					if token == key {

						//value := reflect.Indirect(mapIter.Value())
						value := mapIter.Value()
						temp = value.Interface()
						m = true
						break
					}
				}

				if !m {

					temp = nil
				}

				break
			}
		}

		i++
	}

	return temp
}

// hack valueOf from `reflect.Value` to get indirect as value
// read as interface ptr, ptr, and passing
// better than `reflect.Indirect`

func KIndirectValueOf(data reflect.Value) reflect.Value {

	switch data.Kind() {

	case reflect.Interface:

		// recursion if contain interface, or ptr again
		//return KIndirectValueOf(data.Elem())

		// safety
		data = data.Elem()

		if data.Kind() == reflect.Ptr {

			return data.Elem()
		}

		return data

	case reflect.Ptr:

		return data.Elem()
	}

	return data
}

func KMapSetValue(name string, data any, mapping KMap) bool {

	tokens := strings.Split(name, ".")

	temp := reflect.ValueOf(mapping)

	i := 0
	n := len(tokens)

	if n == 0 {

		return false
	}

	// set value as pointer ? im don't thing so, because can set it without catch pointer
	// `array`, `slice`, or `map` that combine ptr on inside memory
	// that get previous and set value

	for {

		if temp.Interface() == nil {

			return false
		}

		if n <= i {
			break
		}

		token := tokens[i]

		// that have problem, type hide on inside interface and ptr
		// how to solve, is catch all elems on inside interface and set back in var `temp`
		// reset type in var `t`
		temp = KIndirectValueOf(temp)

		// get type
		t := temp.Type()

		// lookup data on `map`
		switch temp.Kind() {

		case reflect.Array, reflect.Slice:

			if index, err := strconv.Atoi(token); err == nil {

				if index < temp.Len() {

					// get previous, set value
					if i+1 == n {

						// get previous, to set value on inside `array` or `slice`
						temp.Index(index).Set(reflect.ValueOf(data))
						return true
					}

					value := temp.Index(index)
					temp = value
					break
				}

				// index out of bound
				return false
			}

			// failed parsing number
			return false

		case reflect.Map:
			if t.Key().Kind() == reflect.String {

				m := false

				mapIter := temp.MapRange()

				for mapIter.Next() {

					key := mapIter.Key()

					if token == key.String() {

						// get previous, set value
						if i+1 == n {

							// get previous, to set value on inside `map`
							temp.SetMapIndex(key, reflect.ValueOf(data))
							return true
						}

						value := mapIter.Value()
						temp = value

						m = true
						break
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

func KMapDelValue(name string, mapping KMap) bool {

	tokens := strings.Split(name, ".")

	var temp, prev reflect.Value

	temp = reflect.ValueOf(mapping)
	prev = reflect.Value{}

	var prevToken string

	i := 0
	n := len(tokens)

	if n == 0 {

		return false
	}

	// set value as pointer ? im don't thing so, because can set it without catch pointer
	// `array`, `slice`, or `map` that combine ptr on inside memory
	// that get previous and set value

	for {

		if temp.Interface() == nil {

			return false
		}

		if n <= i {
			break
		}

		token := tokens[i]

		// that have problem, type hide on inside interface and ptr
		// how to solve, is catch all elems on inside interface and set back in var `temp`
		// reset type in var `t`
		temp = KIndirectValueOf(temp)

		// get type
		t := temp.Type()

		// lookup data on `map`
		switch temp.Kind() {

		case reflect.Array, reflect.Slice:

			if index, err := strconv.Atoi(token); err == nil {

				if index < temp.Len() {

					// get previous, set value
					if i+1 == n {

						// get previous, to set value on inside `array` or `slice`
						L := temp.Slice(0, index)
						R := temp.Slice(panda.Min(index+1, n-1), n)

						k := L.Len() + R.Len()

						// ERROR: problem
						data := reflect.MakeSlice(reflect.SliceOf(prev.Type()), k, k)

						// merging

						//for j := 0; j < L.Len(); j++ {
						//
						//	data.Index(j).Set(L.Index(j))
						//}
						//
						//for j := panda.Min(L.Len()+1, k-1); j < k; j++ {
						//
						//	data.Index(j).Set(L.Index(j))
						//}

						// end

						switch prev.Kind() {

						case reflect.Array, reflect.Slice:

							if mIndex, mErr := strconv.Atoi(prevToken); mErr == nil {

								prev.Index(mIndex).Set(data)
								return true
							}

							return false

						case reflect.Map:
							prev.SetMapIndex(reflect.ValueOf(prevToken), data)
							break
						}

						return true
					}

					value := temp.Index(index)
					prev = temp
					temp = value

					prevToken = token

					break
				}

				// index out of bound
				return false
			}

			// failed parsing number
			return false

		case reflect.Map:
			if t.Key().Kind() == reflect.String {

				m := false

				mapIter := temp.MapRange()

				for mapIter.Next() {

					key := mapIter.Key()

					if token == key.String() {

						// get previous, set value
						if i+1 == n {

							// get previous, to delete value on inside `map`
							temp.SetMapIndex(key, reflect.Value{}) // that am don't know how it works
							return true
						}

						value := mapIter.Value()
						prev = temp
						temp = value

						m = true
						break
					}
				}

				prevToken = token

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
