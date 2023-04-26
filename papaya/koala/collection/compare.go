package collection

import (
  "reflect"
  "skfw/papaya/koala/pp"
  "skfw/papaya/panda"
)

const (
  Unknown = iota
  IsEqual
  IsGreaterThan
  IsLessThan
)

type Compare[T any] struct {
  vA reflect.Value
  vB reflect.Value
}

type CompareImpl[T any] interface {
  Init(a, b T)
  Gt() bool
  Ge() bool
  Lt() bool
  Le() bool
  Eq() bool
}

func CompareNew[T any](a, b T) CompareImpl[T] {

  compare := &Compare[T]{}
  compare.Init(a, b)
  return compare
}

func (c *Compare[T]) Init(a, b T) {

  c.vA = pp.KIndirectValueOf(a)
  c.vB = pp.KIndirectValueOf(b)
}

func (c *Compare[T]) Gt() bool {

  switch CompareValue(c.vA, c.vB) {
  case IsGreaterThan:

    return true
  }

  return false
}

func (c *Compare[T]) Ge() bool {

  switch CompareValue(c.vA, c.vB) {
  case IsGreaterThan, IsEqual:
    return true
  }
  return false
}

func (c *Compare[T]) Lt() bool {

  switch CompareValue(c.vA, c.vB) {
  case IsLessThan:

    return true
  }

  return false
}

func (c *Compare[T]) Le() bool {

  switch CompareValue(c.vA, c.vB) {
  case IsLessThan, IsEqual:

    return true
  }

  return false
}

func (c *Compare[T]) Eq() bool {

  switch CompareValue(c.vA, c.vB) {
  case IsEqual:

    return true
  }

  return false
}

func CompareString(a string, b string) int {

  var i, k, m, n int
  var p, q rune

  m, n = len(a), len(b)
  k = panda.Min(m, n)

  // by value

  for i = 0; i < k; i++ {

    p, q = rune(a[i]), rune(b[i])

    if p > q {

      return IsGreaterThan
    }

    if p < q {

      return IsLessThan
    }
  }

  // by length

  if m > n {

    return IsGreaterThan
  }

  if m < n {

    return IsLessThan
  }

  return IsEqual
}

func CompareValue(vA reflect.Value, vB reflect.Value) int {

  if vA.IsValid() && vB.IsValid() {

    tA := vA.Type()
    tB := vB.Type()

    // same data type
    if tA == tB {

      switch tA.Kind() {
      case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:

        if vA.Int() == vB.Int() {

          return IsEqual
        }

        if vA.Int() > vB.Int() {

          return IsGreaterThan
        }

        if vA.Int() < vB.Int() {

          return IsLessThan
        }

      case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:

        if vA.Uint() == vB.Uint() {

          return IsEqual
        }

        if vA.Uint() > vB.Uint() {

          return IsGreaterThan
        }

        if vA.Uint() < vB.Uint() {

          return IsLessThan
        }

      case reflect.String:

        return CompareString(vA.String(), vB.String())
      }
    }
  }

  return Unknown
}
