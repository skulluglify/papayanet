package main

import (
  "PapayaNet/papaya/koala"
  "PapayaNet/papaya/koala/collection"
)

func main() {

  console := koala.KConsoleNew()

  array := collection.KMiddleListNew[int]()

  array.Add(1, 2, 3, 4, 5, 6, 7, 8, 9, 0)
  //array.Push()
  //array.PushLeft()
  //
  //array.Pop()
  //array.PopLeft()
  //array.Del()
  //array.Splice()
  //
  //array.Slice()

  var i uint

  for i = 0; i < array.Len(); i++ {

    v, e := array.Get(i)
    if e != nil {

      console.Error(e)
      break
    }

    console.Log(i, v)
  }

  array.Push(12)
  array.PushLeft(24)

  console.Log(array.Get(0))
  console.Log(array.Get(array.Len() - 1))

  if m, err := array.Slice(1, 10); err == nil {
    if err = m.ForEach(func(i uint, value int) error {

      console.Log(i, value)

      return nil

    }); err != nil {

      console.Error(err)
    }
  }

  console.Log(array.Get(0))
  console.Log(array.Get(1))

  if m := array.Splice(2, 10, 44, 55); m != nil {
    if err := m.ForEach(func(i uint, value int) error {

      console.Log(i, value, 2)

      return nil

    }); err != nil {

      console.Error(err)
    }
  }

  console.Log(array.Get(0))
  console.Log(array.Get(1))

  console.Log(array.Pop())
  console.Log(array.PopLeft())

  if err := array.ForEach(func(i uint, value int) error {

    console.Log(i, value)

    return nil

  }); err != nil {

    console.Error(err)
  }
}
