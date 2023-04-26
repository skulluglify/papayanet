package main

import (
  "PapayaNet/papaya/koala"
  "PapayaNet/papaya/koala/collection"
  "math"
  "time"
)

func main() {

  a := time.Now().UnixMicro()

  console := koala.KConsoleNew()

  // array := collection.KListNew[int]()
  array := collection.KMiddleListNew[int]()

  n := 6

  for p := 0; p < 20; p++ {

    for i := 0; i < n; i++ {

      k := int(math.Pow(10, float64(i)))

      for j := 0; j < k; j++ {

        array.Push(j)
      }

      for j := 0; j < k; j++ {

        // void
        func(_ int, _ error) {}(array.Get(uint(j)))
      }
    }

    b := time.Now().UnixMicro()

    t := time.UnixMicro(b - a).UTC()

    console.Log(t.Hour(), t.Minute(), t.Second(), t.Nanosecond())
  }
}
