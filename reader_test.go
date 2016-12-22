package swf

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

func TestNewReader(t *testing.T) {
	if reader := NewReader(bytes.NewReader([]byte{1, 2, 3})); reader == nil {
		t.Error("expected non-nil, got nil")
	}
}

func TestReadBits(t *testing.T) {
	reader := NewReader(bytes.NewReader([]byte{0x5f}))
	first, err := reader.ReadBits(3)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if first != 2 {
		t.Errorf("expected 2, got %v", first)
	}

	second, err := reader.ReadBits(5)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if second != 0x1f {
		t.Errorf("expected 0x1f, got %#x", second)
	}

	if _, err := reader.ReadBits(1); err != io.EOF {
		t.Errorf("expected io.EOF, got %v", err)
	}
}

func TestReadInt8(t *testing.T) {
	reader := NewReader(bytes.NewReader([]byte{0x03}))
	v, err := reader.ReadInt8()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if v != 0x03 {
		t.Errorf("expected 0x03, got %#x", v)
	}

	v, err = reader.ReadInt8()
	if err != io.EOF {
		t.Errorf("expected io.EOF, got %v", err)
	}
}

func TestReadInt16(t *testing.T) {
	reader := NewReader(bytes.NewReader([]byte{0x03, 0x72, 0x12}))
	v, err := reader.ReadInt16()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if v != 0x7203 {
		t.Errorf("expected 0x7203, got %#x", v)
	}

	v, err = reader.ReadInt16()
	if err != io.ErrUnexpectedEOF {
		t.Errorf("expected io.ErrUnexpectedEOF, got %v", err)
	}
}

func TestReadInt32(t *testing.T) {
	reader := NewReader(bytes.NewReader([]byte{0x03, 0x72, 0x12, 0x04, 0x12}))
	v, err := reader.ReadInt32()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if v != 0x04127203 {
		t.Errorf("expected 0x04127203, got %#x", v)
	}

	v, err = reader.ReadInt32()
	if err != io.ErrUnexpectedEOF {
		t.Errorf("expected io.ErrUnexpectedEOF, got %v", err)
	}
}

func TestReadUInt8(t *testing.T) {
	reader := NewReader(bytes.NewReader([]byte{0x03}))
	v, err := reader.ReadUInt8()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if v != 0x03 {
		t.Errorf("expected 0x03, got %#x", v)
	}

	v, err = reader.ReadUInt8()
	if err != io.EOF {
		t.Errorf("expected io.EOF, got %v", err)
	}
}

func TestReadUInt16(t *testing.T) {
	reader := NewReader(bytes.NewReader([]byte{0x03, 0x72, 0x12}))
	v, err := reader.ReadUInt16()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if v != 0x7203 {
		t.Errorf("expected 0x7203, got %#x", v)
	}

	v, err = reader.ReadUInt16()
	if err != io.ErrUnexpectedEOF {
		t.Errorf("expected io.ErrUnexpectedEOF, got %v", err)
	}
}

func TestReadUInt32(t *testing.T) {
	reader := NewReader(bytes.NewReader([]byte{0x03, 0x72, 0x12, 0x04, 0x12}))
	v, err := reader.ReadUInt32()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if v != 0x04127203 {
		t.Errorf("expected 0x04127203, got %#x", v)
	}

	v, err = reader.ReadUInt32()
	if err != io.ErrUnexpectedEOF {
		t.Errorf("expected io.ErrUnexpectedEOF, got %v", err)
	}
}

func TestReadEUInt32(t *testing.T) {
	reader := NewReader(bytes.NewReader([]byte{0x5f, 0x8a, 0x89, 0x01, 0x8F}))
	v, err := reader.ReadEUInt32()
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}
	if v != 0x5f {
		t.Errorf("expected 0x5f, got %#x", v)
	}

	v, err = reader.ReadEUInt32()
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}
	if v != 0x448a {
		t.Errorf("expected 0x448a, got %#x", v)
	}

	_, err = reader.ReadEUInt32()
	if err != io.ErrUnexpectedEOF {
		t.Errorf("expected io.ErrUnexpectedEOF, got %v", err)
	}

	_, err = reader.ReadEUInt32()
	if err != io.EOF {
		t.Errorf("expected io.EOF, got %v", err)
	}
}

func TestReadUBitValue(t *testing.T) {
	reader := NewReader(bytes.NewReader([]byte{0x9a})) // 1001 1010
	_, err := reader.ReadUBitValue(34)
	if err == nil || !strings.Contains(err.Error(), "1-32 bits") {
		t.Errorf("expected containing '1-32 bits', got %v", err)
	}
	_, err = reader.ReadUBitValue(0)
	if err == nil || !strings.Contains(err.Error(), "1-32 bits") {
		t.Errorf("expected containing '1-32 bits', got %v", err)
	}

	v, err := reader.ReadUBitValue(5)
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}

	if v != 19 {
		t.Errorf("expected 19, got %v", v)
	}

	v, err = reader.ReadUBitValue(3)
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}
	if v != 2 {
		t.Errorf("expected 2, got %v", v)
	}

	_, err = reader.ReadUBitValue(1)
	if err != io.EOF {
		t.Errorf("expected io.EOF, got %v", err)
	}
}

func TestReadBitValue(t *testing.T) {
	reader := NewReader(bytes.NewReader([]byte{0x9a})) // 1001 1010
	_, err := reader.ReadBitValue(34)
	if err == nil || !strings.Contains(err.Error(), "1-32 bits") {
		t.Errorf("expected containing '1-32 bits', got %v", err)
	}
	_, err = reader.ReadBitValue(0)
	if err == nil || !strings.Contains(err.Error(), "1-32 bits") {
		t.Errorf("expected containing '1-32 bits', got %v", err)
	}

	v, err := reader.ReadBitValue(5)
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}

	if v != -13 {
		t.Errorf("expected -13, got %v", v)
	}

	v, err = reader.ReadBitValue(3)
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}
	if v != 2 {
		t.Errorf("expected 2, got %v", v)
	}

	_, err = reader.ReadUBitValue(1)
	if err != io.EOF {
		t.Errorf("expected io.EOF, got %v", err)
	}
}

func TestReadFixed(t *testing.T) {
	reader := NewReader(bytes.NewReader([]byte{0x00, 0x80, 0x07, 0x00, 0x07, 0x00}))

	v, err := reader.ReadFixed()
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}
	if v != 7.5 {
		t.Errorf("expected 7.5, got %v", v)
	}

	_, err = reader.ReadFixed()
	if err != io.ErrUnexpectedEOF {
		t.Errorf("expected io.ErrUnexpectedEOF, got %v", err)
	}
	_, err = reader.ReadFixed()
	if err != io.EOF {
		t.Errorf("expected io.EOF, got %v", err)
	}
}

func TestReadFixed8(t *testing.T) {
	reader := NewReader(bytes.NewReader([]byte{0x80, 0x09, 0x80, 0xf7, 0x09}))

	v, err := reader.ReadFixed8()
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}
	if v != 9.5 {
		t.Errorf("expected 9.5, got %v", v)
	}

	v, err = reader.ReadFixed8()
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}
	if v != -9.5 {
		t.Errorf("expected 9.5, got %v", v)
	}

	_, err = reader.ReadFixed8()
	if err != io.ErrUnexpectedEOF {
		t.Errorf("expected io.ErrUnexpectedEOF, got %v", err)
	}
	_, err = reader.ReadFixed8()
	if err != io.EOF {
		t.Errorf("expected io.EOF, got %v", err)
	}
}

func TestReadString(t *testing.T) {
	reader := NewReader(bytes.NewReader([]byte{'A', 'B', 'C', 0x0, 'E'}))

	v, err := reader.ReadString()
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}
	if v != "ABC" {
		t.Errorf("expected 'ABC', got %v", v)
	}

	if _, err = reader.ReadString(); err != io.ErrUnexpectedEOF {
		t.Errorf("expected io.ErrUnexpectedEOF, got %v", err)
	}
	if _, err = reader.ReadString(); err != io.EOF {
		t.Errorf("expected io.EOF, got %v", err)
	}
}
