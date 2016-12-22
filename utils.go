package swf

import (
	"io"
)

type simpleByteReadSeeker struct {
	io.ReadSeeker
}

func (r *simpleByteReadSeeker) ReadByte() (byte, error) {
	var b [1]byte
	n, err := r.Read(b[:1])
	if err != nil {
		return 0, err
	}
	if n != 1 {
		return 0, io.ErrUnexpectedEOF
	}
	return b[0], nil
}
