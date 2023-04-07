package types

// [properties] Array

type JSArrayImpl interface {
  Length() int
  Name() string
  Prototype() any
  IsArray(obj any) any
  From(obj any) any    // conversion into new instance of `array`
  Of(elems ...any) any // return `JSArray`
}

// [properties] Array.prototype

type JSArrayPrototypeImpl interface {
  Length()
  Constructor(args ...any) any
  Concat()
  CopyWithin()
  Fill()
  Find()
  FindIndex()
  LastIndexOf()
  Pop()
  Push()
  Reverse()
  Shift()
  Unshift()
  Slice()
  Sort()
  Splice()
  Includes()
  IndexOf()
  Join()
  Keys()
  Entries()
  Values()
  ForEach()
  Filter()
  Flat()
  FlatMap()
  Map()
  Every()
  Some()
  Reduce()
  ReduceRight()
  ToLocaleString()
  ToString()
  At()
  FindLast()
  FindLastIndex()
}
