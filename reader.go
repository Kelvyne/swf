package swf

import (
	"encoding/binary"
	"errors"
	"io"

	"github.com/icza/bitio"
)

// Reader is the minimal interface required to read a swf
type Reader interface {
	io.Reader
	io.Seeker
	ReadByte() (byte, error)
	ReadBits(n uint) (uint32, error)
	ReadInt8() (int8, error)
	ReadInt16() (int16, error)
	ReadInt32() (int32, error)
	ReadUInt8() (uint8, error)
	ReadUInt16() (uint16, error)
	ReadUInt32() (uint32, error)
	ReadEUInt32() (uint32, error)
	ReadBitValue(n uint8) (int32, error)
	ReadUBitValue(n uint8) (uint32, error)
	ReadFixed() (float32, error)
	ReadFixed8() (float32, error)
	ReadString() (string, error)
}

type byteReadSeeker interface {
	io.ReadSeeker
	io.ByteReader
}

type reader struct {
	bitio.Reader
	src byteReadSeeker
}

// NewReader provides a simple way to create a Reader from a given io.Reader
func NewReader(r io.ReadSeeker) Reader {
	if cast, ok := r.(byteReadSeeker); ok {
		return &reader{bitio.NewReader(cast), cast}
	}
	return &reader{bitio.NewReader(r), &simpleByteReadSeeker{r}}
}

func (r *reader) Seek(offset int64, whence int) (int64, error) {
	r.Reader.Align()
	return r.src.Seek(offset, whence)
}

// ReadBits reads n bits
func (r *reader) ReadBits(n uint) (uint32, error) {
	v, err := r.Reader.ReadBits(byte(n))
	return uint32(v), err
}

func (r *reader) read(d interface{}) error {
	r.Align()
	return binary.Read(r, binary.LittleEndian, d)
}

// ReadByte reads a single byte from a io.Reader
func (r *reader) ReadInt8() (int8, error) {
	var v int8
	if err := r.read(&v); err != nil {
		return 0, err
	}
	return v, nil
}

// ReadInt16 reads a signed int 16 from a io.Reader
func (r *reader) ReadInt16() (int16, error) {
	var v int16
	if err := r.read(&v); err != nil {
		return 0, err
	}
	return v, nil
}

// ReadInt32 reads a signed int 32 from a io.Reader
func (r *reader) ReadInt32() (int32, error) {
	var v int32
	if err := r.read(&v); err != nil {
		return 0, err
	}
	return v, nil
}

// ReadUInt8 reads a single unsigned byte from a io.Reader
func (r *reader) ReadUInt8() (uint8, error) {
	var v uint8
	if err := r.read(&v); err != nil {
		return 0, err
	}
	return v, nil
}

// ReadUInt16 reads a signed int 16 from a io.Reader
func (r *reader) ReadUInt16() (uint16, error) {
	var v uint16
	if err := r.read(&v); err != nil {
		return 0, err
	}
	return v, nil
}

// ReadUInt32 reads a signed int 32 from a io.Reader
func (r *reader) ReadUInt32() (uint32, error) {
	var v uint32
	if err := r.read(&v); err != nil {
		return 0, err
	}
	return v, nil
}

// ReadEUInt32 reads a swf encoded unsigned int 32 from a io.Reader
// It reads one more byte while the most significant bit is 1
func (r *reader) ReadEUInt32() (uint32, error) {
	var v uint32
	var count uint32
	for {
		b, err := r.ReadUInt8()
		if err != nil {
			if err == io.EOF && count != 0 {
				err = io.ErrUnexpectedEOF
			}
			return 0, err
		}

		leasts := uint32(b) & 0x7f
		v = v | (leasts << (count * 7))

		count++
		if count == 4 || (b&0x80) == 0 {
			break
		}
	}
	return v, nil
}

// ReadUBitValue reads a swf encoded unsigned bit value with n bits
// from a swf.Reader
func (r *reader) ReadUBitValue(n uint8) (uint32, error) {
	if n > 32 || n == 0 {
		return 0, errors.New("bit value is 1-32 bits")
	}
	return r.ReadBits(uint(n))
}

// ReadBitValue reads a swf encoded signed bit value with n bits
// from a swf.Reader
func (r *reader) ReadBitValue(n uint8) (int32, error) {
	value, err := r.ReadUBitValue(n)
	if err != nil {
		return 0, err
	}

	sign := (value >> (n - 1)) & 0x1
	if sign == 1 {
		value = value | (0xffffffff << n)
	}

	return int32(value), nil
}

// ReadFixed reads a swf encoded fixed point number from a io.Reader
// Each part of the fixed point number is 16 bits
func (r *reader) ReadFixed() (float32, error) {
	// float16 are not in Go, so read uint16 (little endian) then cast it
	after, err := r.ReadUInt16()
	if err != nil {
		return 0, err
	}

	before, err := r.ReadUInt16()
	if err != nil {
		if err == io.EOF {
			err = io.ErrUnexpectedEOF
		}
		return 0, err
	}

	return float32(before) + float32(after)/65536, nil
}

// ReadFixed8 reads a swf encoded fixed point number from a io.Reader.
// Each part of the fixed point number is 8 bits
func (r *reader) ReadFixed8() (float32, error) {
	after, err := r.ReadUInt8()
	if err != nil {
		return 0, err
	}
	before, err := r.ReadInt8()
	if err != nil {
		if err == io.EOF {
			return 0, io.ErrUnexpectedEOF
		}
	}
	if before < 0 {
		return float32(before) - float32(after)/256, nil
	}
	return float32(before) + float32(after)/256, nil
}

func (r *reader) ReadString() (string, error) {
	var bytes []byte
	for {
		c, err := r.ReadUInt8()
		if err != nil {
			if err == io.EOF && len(bytes) > 0 {
				return "", io.ErrUnexpectedEOF
			}
			return "", err
		}
		if c == 0x0 {
			break
		}
		bytes = append(bytes, byte(c))
	}

	return string(bytes), nil
}
