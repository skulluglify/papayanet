package main

import (
  "PapayaNet/papaya/koala"
  "PapayaNet/papaya/koala/collection"
  "PapayaNet/papaya/koala/mapping"
  "math"
)

func main() {

  var i uint
  console := koala.KConsoleNew()
  console.Log("KList testing ...")

  list := collection.KListNewR[int]([]int{12, 24, 36, 48, 60, 72, 84, 96, 108})

  for i = 0; i < list.Len(); i++ {

    console.Log(list.Get(i))
  }

  console.Warn("splice ...")

  // TODO: fix problem --
  // TODO: initial val.Type() after val.IsValid() ++
  removes := list.Splice(1, 1, 72, 80)

  for i = 0; i < removes.Len(); i++ {

    console.Log(removes.Get(i))
  }

  console.Warn("look ...")
  console.Warn(list.Len())

  for i = 0; i < list.Len(); i++ {

    console.Log(list.Get(i))
  }

  type Group struct {
    Message string
  }

  data := &mapping.KMap{
    "title":       "Koala Mapping",
    "shortname":   "KMap",
    "description": "Koala Mapping for defined new Object",
    "docs": &mapping.KMap{
      "lesson1": []string{
        "make it possible with pointer or references",
        "deeper collection of enums in map object",
        "like JSON",
      },
    },
    "premises": []mapping.KMap{
      {
        "read":  "wait",
        "write": "wait",
      },
    },
    "group": Group{
      Message: "this is group for map Object",
    },
  }

  for _, enum := range data.Tree().Enums() {

    k, v := enum.Tuple()
    console.Log(k, v)
  }

  for _, k := range data.Tree().Keys() {

    console.Log(k)
    console.Log(data.Get(k))
  }

  console.Log("--- iter ---")

  iter := data.Tree().Iterable()
  for next := iter.Next(); next.HasNext(); next = next.Next() {

    k, v := next.Enum().Tuple()

    console.Log(k, v)
  }

  // path, prefix, suffix
  // Index
  // Get, Set, Del, Put

  console.Log(data.Put("math", math.Pi))
  console.Log(data.Put("premises.0.read", "ready"))
  console.Log(data.Put("premises.0.shared", nil))

  console.Log(data.Put("system.dev.name", "ubuntu"))

  data.Branch("system.dev.freedesktop.x11.cosmos")

  console.Log("--- json ---")
  console.Log(data.JSON())
}
