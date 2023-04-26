package main

import (
  "skfw/papaya/koala"
  "skfw/papaya/koala/collection"
)

func main() {

  console := koala.KConsoleNew()

  array := collection.KMiddleListNew[*[]rune]()

  for _, w := range []string{"ba", "cad", "bab", "aku", "a", "bc", "z", "zo"} {

    y := []rune(w)
    array.Push(&y)
  }
  //array.Add("ba", "cad", "bab", "aku", "a", "bc", "z", "zo")
  array.Sort()
  //array.Reverse()

  if err := array.ForEach(func(i uint, value *[]rune) error {

    if i > 10 {

      return nil
    }

    console.Log(i, string(*value))
    return nil

  }); err != nil {

    return
  }

  var i uint

  for i = 0; i < array.Len(); i++ {

    console.Log(array.Get(i))
  }
}
