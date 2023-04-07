package kio

import (
  "os"
)

type KFile struct {
  Path string
}

type KFileImpl interface {
  Cat() string
  IsExists() bool
}

func KFileNew(path string) KFileImpl {

  file := &KFile{
    Path: path,
  }

  return file
}

func (f *KFile) Cat() string {

  data, _ := os.ReadFile(f.Path)
  return string(data)
}

func (f *KFile) IsExists() bool {

  if _, err := os.Stat(f.Path); err != nil {

    //if os.IsNotExist(err) {
    //
    //  return false
    //}
    ////error checking
    //panic(err)

    return false
  }

  return true
}
