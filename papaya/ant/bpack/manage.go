package bpack

import (
  "errors"
  "os"
  "skfw/papaya/koala/kio"
  "skfw/papaya/koala/tools/posix"
)

var FileNotFound = errors.New("file not found")

func OpenPacket(path string) *Packet {

  for _, packet := range Pkts {

    if packet.Path == path {

      return &packet
    }
  }

  return nil
}

func OpenFile(path string) ([]byte, error) {

  // try to read from buffer
  var data []byte

  if pkt := OpenPacket(path); pkt != nil {

    return pkt.Data, nil
  }

  // fake root directory
  var rootDir string

  if rootDir = FindDataPath(PATH); rootDir != "" {

    path = posix.KPathNew(path).JoinStr(rootDir)
  }

  // check existing
  var file kio.KFileImpl

  file = kio.KFileNew(path)

  if file.IsExist() {

    if file.IsFile() {

      return os.ReadFile(path)
    }
  }

  return data, FileNotFound
}
