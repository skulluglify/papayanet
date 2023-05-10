package pp

// Method for Inline Statement

func L[T any](check bool, good T, bad T) T {

  if check {

    return good
  }

  return bad
}

// auto type defined by name

var Lany = L[any]
var Lstr = L[string]
var Lbool = L[bool]
var Lbyte = L[byte]
var Lint = L[int]
var Luint = L[uint]
