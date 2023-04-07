package pp

// Method for Inline Statement

func KIS[T any](check bool, good T, bad T) T {

  if check {

    return good
  }

  return bad
}

// auto type defined by name

var KISAny = KIS[any]
var KISStr = KIS[string]
var KISBool = KIS[bool]
var KISByte = KIS[byte]
var KISInt = KIS[int]
var KISUint = KIS[uint]
