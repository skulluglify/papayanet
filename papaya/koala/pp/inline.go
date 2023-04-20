package pp

// Method for Inline Statement

func L[T any](check bool, good T, bad T) T {

  if check {

    return good
  }

  return bad
}

// auto type defined by name

var LAny = L[any]
var LStr = L[string]
var LBool = L[bool]
var LByte = L[byte]
var LInt = L[int]
var LUint = L[uint]
