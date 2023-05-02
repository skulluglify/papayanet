package bpack

import (
  "errors"
  "skfw/papaya/koala/kio"
  "skfw/papaya/koala/tools/posix"
  "os"
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

  // try read from buffer
  var data []byte

  if pkt := OpenPacket(path); pkt != nil {

    return pkt.Data, nil
  }

  // fake rootdir
  var rootdir string

  if rootdir = FindDataPath(PATH); rootdir != "" {

    path = posix.KPathNew(path).JoinStr(rootdir)
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
