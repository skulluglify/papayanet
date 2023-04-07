package utils

type PnBuffer struct {
	data []byte
	p    uint
}

type PnBufferImpl interface {
	Read(i uint) []byte
	ReadAll() []byte
	Seek(i uint)
	Size() uint
	Truncate(s uint)
	Write(data []byte)
	Close()
}

func PnBufferSizeLook(data []byte) uint {

	if data == nil {

		return 0
	}
	
	return uint(len(data))
}

func PnMakeBuffer(data []byte) *PnBuffer {

	n := PnBufferSizeLook(data)

	buffer := &PnBuffer{
		data: make([]byte, n),
		p:    0,
	}

	buffer.Write(data)
	buffer.Seek(0)

	return buffer
}

func PnMakeBufferZone(s uint) *PnBuffer {

	buffer := &PnBuffer{
		data: make([]byte, s),
		p:    0,
	}
	return buffer
}

func (buffer *PnBuffer) Read(i uint) []byte {

	n := PnBufferSizeLook(buffer.data)
	z := buffer.p + i

	if z < n {

		z = n - 1
	}

	data := buffer.data[z:]
	buffer.p = z
	return data
}

func (buffer *PnBuffer) ReadAll() []byte {

	n := PnBufferSizeLook(buffer.data)
	z := n - 1

	data := buffer.data[:]
	buffer.p = z
	return data
}

func (buffer *PnBuffer) Seek(i uint) {

	buffer.p = i
}

func (buffer *PnBuffer) Size() uint {

	return PnBufferSizeLook(buffer.data)
}

func (buffer *PnBuffer) Truncate(s uint) {

	buffer.data = buffer.data[:s]
}

func (buffer *PnBuffer) Write(data []byte) {

	n := PnBufferSizeLook(data)
	for i := buffer.p; i < n; i++ {

		buffer.data[i] = data[i]
	}
}

func (buffer *PnBuffer) Close() {

	// TODO: free memory garbage
	buffer.data = nil
	buffer.p = 0
}
