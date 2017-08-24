package bytecode

import (
	"encoding/binary"
	"io"
)

// Writer is used to serialize as3 bytecode
type Writer interface {
	WriteU8(uint8) error
	WriteU16(uint16) error
	WriteS24(int32) error
	WriteU30(uint32) error
	WriteU32(uint32) error
	WriteS32(int32) error
	WriteD64(float64) error
}

type writer struct {
	io.Writer
}

// NewWriter provides a simple way to create an as3 bytecode writer
func NewWriter(w io.Writer) Writer {
	return &writer{w}
}

func (w *writer) write(data interface{}) error {
	return binary.Write(w.Writer, binary.LittleEndian, data)
}

func (w *writer) WriteU8(x uint8) error {
	return w.write(x)
}

func (w *writer) WriteU16(x uint16) error {
	return w.write(x)
}

func (w *writer) WriteS24(x int32) error {
	var bytes [3]byte
	bytes[0] = byte(x & 0xff)
	bytes[1] = byte((x >> 8) & 0xff)
	bytes[2] = byte((x >> 16) & 0xff)

	return w.write(bytes)
}

func (w *writer) writeVariableLength(x uint32) error {
	var bytes [5]byte
	var n uint32
	// always write at least 1 byte
	bytes[n] = byte(x & 0x7f)
	n++
	for (x >> (n * 7)) != 0 {
		// set previous byte as not last
		bytes[n-1] |= 0x80
		bytes[n] = byte((x >> (n * 7)) & 0x7f)
		n++
	}
	return w.write(bytes[0:n])
}

func (w *writer) WriteU30(x uint32) error {
	return w.WriteU32(x)
}

func (w *writer) WriteU32(x uint32) error {
	return w.writeVariableLength(x)
}

func (w *writer) WriteS32(x int32) error {
	return w.writeVariableLength(uint32(x))
}

func (w *writer) WriteD64(x float64) error {
	return w.write(x)
}
