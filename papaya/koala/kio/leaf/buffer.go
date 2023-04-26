package leaf

import (
  "reflect"
  "skfw/papaya/panda/nosign"
)

type KBuffer struct {
  data []byte
  p    uint
}

var KBufferType = reflect.TypeOf(KBuffer{})

type KBufferImpl interface {
  Read(i uint) []byte
  ReadAll() []byte
  Seek(i uint)
  Size() uint
  Truncate(s uint)
  Write(data []byte)
  Close()
}

func KBufferSizeLook(data []byte) uint {

  if data == nil {

    return 0
  }

  return uint(len(data))
}

func KMakeBuffer(data []byte) *KBuffer {

  n := KBufferSizeLook(data)

  buffer := &KBuffer{
    data: make([]byte, n),
    p:    0,
  }

  buffer.Write(data)
  buffer.Seek(0)

  return buffer
}

func KMakeBufferZone(s uint) *KBuffer {

  buffer := &KBuffer{
    data: make([]byte, s),
    p:    0,
  }
  return buffer
}

func (b *KBuffer) Read(i uint) []byte {

  n := KBufferSizeLook(b.data)
  z := b.p + i

  if z < n {

    z = n - 1
  }

  data := b.data[z:]
  b.p = z
  return data
}

func (b *KBuffer) ReadAll() []byte {

  n := KBufferSizeLook(b.data)
  z := n - 1

  data := b.data[:]
  b.p = z
  return data
}

func (b *KBuffer) Seek(i uint) {

  b.p = i
}

func (b *KBuffer) Size() uint {

  return KBufferSizeLook(b.data)
}

func (b *KBuffer) Truncate(s uint) {

  b.data = b.data[:s]

  // update point
  if s > 0 {

    s-- // last index
    b.p = nosign.Min(s, b.p)

  } else {

    b.p = 0
  }
}

func (b *KBuffer) Write(data []byte) {

  n := KBufferSizeLook(data)
  for i := b.p; i < n; i++ {

    b.data[i] = data[i]
  }
}

func (b *KBuffer) Close() {

  // TODO: free memory garbage
  b.data = nil
  b.p = 0
}
