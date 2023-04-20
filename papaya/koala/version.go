package koala

import (
  "strconv"
  "strings"
)

type KVersion struct {
  Major uint8
  Minor uint8
  Patch uint8
}

type KVersionImpl interface {
  String() string
}

func KVersionNew(major uint8, minor uint8, patch uint8) KVersionImpl {

  return &KVersion{
    Major: major,
    Minor: minor,
    Patch: patch,
  }
}

func (version *KVersion) String() string {

  return strings.Join([]string{
    strconv.Itoa(int(version.Major)),
    strconv.Itoa(int(version.Minor)),
    strconv.Itoa(int(version.Patch))},
    ".")
}
