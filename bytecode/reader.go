package bytecode

import "io"
import "encoding/binary"
import "errors"

// ErrMalformedVariableInteger means that an encoded variable integer is
// malformed (its length is >5 bytes)
var ErrMalformedVariableInteger = errors.New("malformed variable integer")

// Reader is the minimal interface required to read as3 bytecode
type Reader interface {
	ReadU8() (uint8, error)
	ReadU16() (uint16, error)
	ReadS24() (int32, error)
	ReadU30() (uint32, error)
	ReadU32() (uint32, error)
	ReadS32() (int32, error)
	ReadD64() (float64, error)
	ReadBytes(n uint32) ([]byte, error)
}

type reader struct {
	io.Reader
}

// NewReader provides a simple way to create an as3 bytecode reader
func NewReader(r io.Reader) Reader {
	return &reader{r}
}

func (r *reader) read(data interface{}) error {
	return binary.Read(r.Reader, binary.LittleEndian, data)
}

func (r *reader) ReadU8() (uint8, error) {
	var v uint8
	err := r.read(&v)
	return v, err
}

func (r *reader) ReadU16() (uint16, error) {
	var v uint16
	err := r.read(&v)
	return v, err
}

func (r *reader) ReadS24() (int32, error) {
	var bytes [3]byte
	err := r.read(bytes[:])
	if err != nil {
		return 0, err
	}
	v := uint32(bytes[0]) | uint32(bytes[1])<<8 | uint32(bytes[2])<<16
	if v>>23 != 0 {
		v |= 0xff000000
	}
	return int32(v), nil
}

func (r *reader) readVariableLength() (v uint32, n uint8, err error) {
	var b uint8 = 0x80
	for b>>7 != 0 {
		if n >= 5 {
			err = ErrMalformedVariableInteger
			v = 0
			return
		}
		b, err = r.ReadU8()
		if err != nil {
			return
		}
		v |= (uint32(b & 0x7f)) << (n * 7)
		n++
	}
	return
}

func (r *reader) ReadU32() (uint32, error) {
	v, _, err := r.readVariableLength()
	return v, err
}

func (r *reader) ReadU30() (uint32, error) {
	return r.ReadU32()
}

func (r *reader) ReadS32() (int32, error) {
	v, n, err := r.readVariableLength()
	if err != nil {
		return 0, err
	}
	// If the higher bit of the last read byte is 1, it means we need to
	// expand the sign (v is negative)
	if v&(1<<(n*7-1)) != 0 {
		// shift should be n*7 + 1 but we know n*7th bit is set
		v |= 0xffffffff << (n * 7)
	}
	return int32(v), nil
}

func (r *reader) ReadD64() (float64, error) {
	var v float64
	err := r.read(&v)
	return v, err
}

func (r *reader) ReadBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := io.ReadFull(r.Reader, b)
	if err != nil {
		return nil, err
	}
	return b, nil
}
