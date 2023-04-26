package main

import (
  "skfw/papaya/koala"
  "skfw/papaya/koala/collection"
  "skfw/papaya/koala/mapping"
)

func main() {

  console := koala.KConsoleNew()

  value := collection.ValueNew(&mapping.KMap{
    "message": "hello, world!",
  })

  console.Log(value.Map())
}
