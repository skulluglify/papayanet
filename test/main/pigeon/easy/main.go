package main

import (
  "fmt"
  "skfw/papaya/koala/pp"
  "skfw/papaya/pigeon/easy"
)

type Test struct {
  ID string
}

func (Test) TableName() string {

  return "test"
}

func main() {

  test := &Test{
    ID: "uniq",
  }

  fmt.Println(easy.TableName(test))

  val := pp.KIndirectValueOf(test)
  xID := val.FieldByName("ID")

  fmt.Println(xID)
  fmt.Println(easy.FieldGet(xID))
  easy.FieldSet(xID, "rand")
  fmt.Println(easy.FieldGet(xID))
}
