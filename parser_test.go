package swf

import (
	"bytes"
	"io"
	"os"
	"reflect"
	"testing"
)

const (
	fixturePath = "./fixtures/big.swf"
)

func createReader(t *testing.T) *os.File {
	file, err := os.Open(fixturePath)
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	return file
}

func TestParse(t *testing.T) {
	file := createReader(t)
	defer file.Close()

	p := NewParser(file)

	swf, err := p.Parse()
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}

	correctHeader := Header{
		CompressionZlib,
		11, 11605652,
		Rect{16, 0, 25600, 0, 20480},
		50.0, 1,
	}
	if !reflect.DeepEqual(swf.Header, correctHeader) {
		t.Errorf("expected %v, got %v", correctHeader, swf.Header)
	}
	doAbc, ok := swf.Tags[0].(*TagDoABC)
	if !ok {
		t.Errorf("expected %v to be a *TagDoABC", swf.Tags[0])
	}

	if doAbc.Name != "frame1" {
		t.Errorf("expected 'frame1', got %v", doAbc.Name)
	}
	if len(doAbc.ABCData) != 7973009 {
		t.Errorf("expected 7973009, got %v", len(doAbc.ABCData))
	}
}

func TestParseHeader(t *testing.T) {
	headerBytes := []byte{
		0x46, 0x57, 0x53,
		0x0b,
		0x94, 0x16, 0xb1, 0x00,
		0x80, 0x00, 0x03, 0x20, 0x00, 0x00, 0x02, 0x80, 0x00,
		0x00, 0x32,
		0x01, 0x00,
	}
	p := newParser(bytes.NewReader(headerBytes))

	header, err := p.ParseHeader()
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}

	correctHeader := Header{
		CompressionNone,
		11, 11605652,
		Rect{16, 0, 25600, 0, 20480},
		50.0, 1,
	}
	if !reflect.DeepEqual(header, correctHeader) {
		t.Errorf("expected %v, got %v", correctHeader, header)
	}
}

func TestParseRect(t *testing.T) {
	rectBytes := []byte{0x80, 0x00, 0x03, 0x20, 0x00, 0x00, 0x02, 0x80, 0x00}
	p := newParser(bytes.NewReader(rectBytes))
	rect, err := p.ParseRect()
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}

	if rect.NBits != 16 {
		t.Errorf("expected 16, got %v", rect.NBits)
	}
	if rect.Xmin != 0 {
		t.Errorf("expected 0, got %v", rect.Xmin)
	}
	if rect.Xmax != 25600 {
		t.Errorf("expected 25600, got %v", rect.Xmax)
	}
	if rect.Ymin != 0 {
		t.Errorf("expected 0, got %v", rect.Xmin)
	}
	if rect.Ymax != 20480 {
		t.Errorf("expected 20480, got %v", rect.Xmax)
	}

	var fails = [][]byte{
		rectBytes[0:1],
		rectBytes[0:3],
		rectBytes[0:5],
		rectBytes[0:7],
	}

	for _, fail := range fails {
		func(b []byte) {
			pFail := newParser(bytes.NewReader(b))
			_, err := pFail.ParseRect()
			if err != io.ErrUnexpectedEOF {
				t.Errorf("expected io.ErrUnexpectedEOF, got %v", err)
			}
		}(fail)
	}
}

func TestHandleEOF(t *testing.T) {
	pFail := newParser(bytes.NewReader([]byte{0}))
	if err := pFail.handleEOF(nil); err != nil {
		t.Errorf("expected nil, got %v", err)
	}

	if err := pFail.handleEOF(io.ErrNoProgress); err != io.ErrNoProgress {
		t.Errorf("expected io.ErroNoProgress, got %v", err)
	}

	if err := pFail.handleEOF(io.EOF); err != io.ErrUnexpectedEOF {
		t.Errorf("expected io.ErrUnexpectedEOF, got %v", err)
	}

	if err := pFail.handleEOF(io.ErrUnexpectedEOF); err != io.ErrUnexpectedEOF {
		t.Errorf("expected io.ErrUnexpectedEOF, got %v", err)
	}
}
