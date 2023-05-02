package bpack

type Packet struct {
  Path string
  Mimetype string
  Charset string
  Data []byte
}

type Packets []Packet
