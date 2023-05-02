package bpack

type Packet struct {
  Path     string
  Mimetype string
  Charset  string
  Data     []byte
  Size     uint64
}

type Packets []Packet
