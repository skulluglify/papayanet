package test

import (
  "PapayaNet/papaya/koala"
  "PapayaNet/papaya/koala/gen"
  "PapayaNet/papaya/koala/mapping"
  "testing"
)

func TestMap(t *testing.T) {

  data := mapping.KMap{
    "name":        "Koala",
    "version":     koala.KVersionNew(1, 0, 0),
    "description": "Library Koala, Make it dynamic typing",
    "system": &mapping.KMap{
      "name":      "Windows Selinux",
      "shortname": "WSL",
      "distro":    "Ubuntu",
      "machine":   "amd64",
      "UEFI":      true,
    },
    "dependencies": mapping.KMap{
      "panda": mapping.KMap{
        "versions": []string{
          "v1.0.0",
          "v2.0.0",
          "v3.0.0",
        },
      },
    },
  }

  t.Log("test catch values")

  if mapping.KMapGetValue("name", data) != "Koala" {

    t.Error("key `name` wrong value")
  }

  if mapping.KMapGetValue("system.name", data) != "Windows Selinux" {

    t.Error("key `system.name` wrong value")
  }

  if mapping.KMapGetValue("dependencies.panda.versions.0", data) != "v1.0.0" {

    t.Error("key `dependencies.panda.versions.0` wrong value")
  }

  t.Log("test change values")

  mapping.KMapSetValue("name", "koala", data)
  mapping.KMapSetValue("system.name", "Windows 11", data)
  mapping.KMapSetValue("dependencies.panda.versions.0", "v1.2.0", data)

  if mapping.KMapGetValue("name", data) != "koala" {

    t.Error("key `name` wrong value")
  }

  if mapping.KMapGetValue("system.name", data) != "Windows 11" {

    t.Error("key `system.name` wrong value")
  }

  if mapping.KMapGetValue("dependencies.panda.versions.0", data) != "v1.2.0" {

    t.Error("key `dependencies.panda.versions.0` wrong value")
  }

  t.Log("try delete key `system.shortname`", mapping.KMapDelValue("system.shortname", data))
  t.Log("try delete key `dependencies.panda.versions.1`", mapping.KMapDelValue("dependencies.panda.versions.1", data))
  t.Log("try delete key `dependencies.panda.versions.1`", mapping.KMapDelValue("dependencies.panda.versions.1", data))
  t.Log("try delete key `dependencies.panda.versions.2`", mapping.KMapDelValue("dependencies.panda.versions.2", data))

  t.Log(mapping.KMapKeys(data))
  t.Log(mapping.KMapTreeKeys(data))

  if mapping.KMapGetValue("system.shortname", data) != nil {

    t.Error("key `system.shortname` cannot be removed")
  }

  if mapping.KMapGetValue("dependencies.panda.versions.1", data) != nil {

    t.Error("key `dependencies.panda.versions.1` cannot be removed")
  }

  if mapping.KMapGetValue("dependencies.panda.versions.2", data) != nil {

    t.Error("key `dependencies.panda.versions.2` cannot be removed")
  }

  iter := gen.KMapTreeIterable(&data)

  for next := iter.Next(); next.HasNext(); next = next.Next() {

    t.Log(next.Enum())
  }

  t.Log(mapping.KMapValues(data))
}
