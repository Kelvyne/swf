package swf

import (
	"bytes"
	"compress/zlib"
	"errors"
	"io"
	"io/ioutil"
)

// ErrMalformedHeader means that the swf file header is malformed.
// The signature is malformed
var ErrMalformedHeader = errors.New("malformed header")

// ErrUnsupportedFile means that the swf file is not supported.
// The file is compressed with an unsupported algorithm
var ErrUnsupportedFile = errors.New("unsupported file")

// Parser is the minimal interface for parsing a Swf file
type Parser interface {
	Parse() (Swf, error)
}

type parser struct {
	r      Reader
	origin io.ReadSeeker
}

func newParser(origin io.ReadSeeker) *parser {
	return &parser{NewReader(origin), origin}
}

// Parse creates a Parser and parses the given input
func Parse(origin io.ReadSeeker) (Swf, error) {
	return newParser(origin).Parse()
}

// NewParser provides a simple way to create a Swf file
func NewParser(origin io.ReadSeeker) Parser {
	return newParser(origin)
}

func (p *parser) handleEOF(err error) error {
	if err == io.EOF {
		return io.ErrUnexpectedEOF
	}
	return err
}

// Parse parses an entire Swf file
func (p *parser) Parse() (s Swf, err error) {
	header, err := p.ParseHeader()
	if err != nil {
		return
	}
	s.Header = header
	if s.Header.Compression != CompressionNone {
		defer func() {
			if err == nil {
				err = p.origin.(io.ReadCloser).Close()
			}
		}()
	}

	tags, err := p.ParseTags()
	if err != nil {
		return
	}
	s.Tags = tags

	return
}

func (p *parser) replaceReader(compression uint8) error {
	switch compression {
	default:
		break
	case CompressionZlib:
		_, err := p.origin.Seek(8, io.SeekStart)
		if err != nil {
			return err
		}
		r, err := zlib.NewReader(p.origin)
		if err != nil {
			return err
		}
		buf, err := ioutil.ReadAll(r)
		if err != nil {
			return err
		}
		if err = r.Close(); err != nil {
			return err
		}
		p.r = NewReader(bytes.NewReader(buf))
	}
	return nil
}

func (p *parser) ParseHeader() (Header, error) {
	signature, err := p.r.ReadUInt8()
	var compression uint8

	if err != nil {
		return Header{}, p.handleEOF(err)
	}
	switch signature {
	default:
		return Header{}, ErrMalformedHeader
	case 'F':
		compression = CompressionNone
	case 'C':
		compression = CompressionZlib
	case 'Z':
		return Header{}, ErrUnsupportedFile
	}

	if signature, err = p.r.ReadUInt8(); err != nil {
		return Header{}, p.handleEOF(err)
	} else if signature != 'W' {
		return Header{}, ErrMalformedHeader
	}
	if signature, err = p.r.ReadUInt8(); err != nil {
		return Header{}, p.handleEOF(err)
	} else if signature != 'S' {
		return Header{}, ErrMalformedHeader
	}
	version, err := p.r.ReadUInt8()
	if err != nil {
		return Header{}, err
	}
	fileLength, err := p.r.ReadUInt32()
	if err != nil {
		return Header{}, p.handleEOF(err)
	}

	if err = p.replaceReader(compression); err != nil {
		return Header{}, p.handleEOF(err)
	}

	frameSize, err := p.ParseRect()
	if err != nil {
		return Header{}, p.handleEOF(err)
	}

	frameRate, err := p.r.ReadFixed8()
	if err != nil {
		return Header{}, p.handleEOF(err)
	}
	frameCount, err := p.r.ReadUInt16()
	if err != nil {
		return Header{}, p.handleEOF(err)
	}
	return Header{compression, version, fileLength, frameSize, frameRate, frameCount}, nil
}

func (p *parser) ParseTags() ([]Tag, error) {
	var tags []Tag

	for {
		t, err := p.ParseTag()
		if err != nil {
			return nil, err
		}
		if t != nil {
			tags = append(tags, t)
			if t.Code() == CodeTagEnd {
				break
			}
		}
	}

	return tags, nil
}

func (p *parser) ParseTag() (Tag, error) {
	codeAndLength, err := p.r.ReadUInt16()
	if err != nil {
		return nil, p.handleEOF(err)
	}
	code := (codeAndLength >> 6) & 0x3ff
	length := uint32(codeAndLength & 0x3f)
	if length == 0x3f {
		length, err = p.r.ReadUInt32()
		if err != nil {
			return nil, p.handleEOF(err)
		}
	}

	type handleFunc func(uint32) (Tag, error)
	supportedTags := map[uint16]handleFunc{
		CodeTagEnd:   p.ParseTagEnd,
		CodeTagDoABC: p.ParseTagDoABC,
	}

	if handler, found := supportedTags[code]; found {
		t, err := handler(length)
		if err != nil {
			return nil, p.handleEOF(err)
		}
		return t, nil
	}

	// Discard data since we can not handle it
	if n, err := io.CopyN(ioutil.Discard, p.r, int64(length)); err != nil {
		return nil, p.handleEOF(err)
	} else if uint32(n) != length {
		return nil, io.ErrUnexpectedEOF
	}
	return nil, nil
}

func (p *parser) ParseTagEnd(length uint32) (Tag, error) {
	return &tag{CodeTagEnd, length}, nil
}

func (p *parser) ParseTagDoABC(length uint32) (Tag, error) {
	begin, err := p.r.Seek(0, io.SeekCurrent)
	if err != nil {
		return nil, err
	}
	flags, err := p.r.ReadUInt32()
	if err != nil {
		return nil, p.handleEOF(err)
	}
	name, err := p.r.ReadString()
	if err != nil {
		return nil, p.handleEOF(err)
	}
	end, err := p.r.Seek(0, io.SeekCurrent)
	if err != nil {
		return nil, err
	}
	abcDataLen := length - uint32(end-begin)
	abcData := make([]byte, abcDataLen)
	if n, err := io.ReadFull(p.r, abcData); err != nil {
		return nil, p.handleEOF(err)
	} else if uint32(n) != abcDataLen {
		return nil, io.ErrUnexpectedEOF
	}
	return &TagDoABC{tag{CodeTagDoABC, length}, flags, name, abcData}, nil
}

// ParseRect parses a Rectangle record
func (p *parser) ParseRect() (rect Rect, err error) {
	nBits, err := p.r.ReadUBitValue(5)
	if err != nil {
		return
	}
	rect.NBits = uint8(nBits)

	readField := func(ptr *int32) error {
		value, fieldErr := p.r.ReadBitValue(rect.NBits)
		if fieldErr == io.EOF {
			fieldErr = io.ErrUnexpectedEOF
		}
		*ptr = value
		return fieldErr
	}

	if err = readField(&rect.Xmin); err != nil {
		return
	}
	if err = readField(&rect.Xmax); err != nil {
		return
	}
	if err = readField(&rect.Ymin); err != nil {
		return
	}
	if err = readField(&rect.Ymax); err != nil {
		return
	}
	return
}
