package swf

// These represent code of handled Swf tags
const (
	// CodeTagEnd is the code representing a Tag of type End
	CodeTagEnd = 0
	// CodeTagDoABC is the code representing a Tag of type DoABC
	CodeTagDoABC = 82
)

// These represent possible Swf compression algorithm
const (
	CompressionNone = iota
	CompressionZlib
	CompressionLZMA
)

// Swf represents a Swf file deserialized
type Swf struct {
	Header Header
	Tags   []Tag
}

// Header represents a Swf file's header
type Header struct {
	Compression uint8
	Version     uint8
	FileLength  uint32
	FrameSize   Rect
	FrameRate   float32
	FrameCount  uint16
}

// Tag represents the generic interface for representing a Swf Tag
type Tag interface {
	Code() uint16
	Length() uint32
}

type tag struct {
	code   uint16
	length uint32
}

// TagDoABC represents a DoABC Tag
type TagDoABC struct {
	tag
	Flags   uint32
	Name    string
	ABCData []byte
}

func (t *tag) Code() uint16   { return t.code }
func (t *tag) Length() uint32 { return t.length }

// Rect represents a Rectangle record
type Rect struct {
	NBits uint8
	Xmin  int32
	Xmax  int32
	Ymin  int32
	Ymax  int32
}
