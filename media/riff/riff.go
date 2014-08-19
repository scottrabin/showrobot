package riff

import (
	"bytes"
	"fmt"
	"os"
)

// A RIFF element is one of the three components comprising a RIFF file:
// a RIFF header, a LIST, or a CHUNK
type riffElement interface {
	fmt.Stringer
}

// A FOURCC (four-character code) is a 32-bit unsigned integer created
// by concatenating four ASCII characters
type FOURCC [4]byte

// A RIFF header has the form `'RIFF' fileSize fileType data`, where:
// 'RIFF' is the literal FOURCC 'RIFF' (omitted from the struct),
// fileSize is a 4-byte value representing the size of the file (excluding
//          the 'RIFF' FOURCC or the 4 bytes of fileSize),
// fileType is a FOURCC that identifies the specific type of file
type riffheader struct {
	fileSize uint32
	fileType FOURCC
}

func (h *riffheader) String() string {
	return fmt.Sprintf("RIFF '%s' (%d)", string(h.fileType[:]), h.fileSize)
}

// A chunk has the form `ckID ckSize ckData`, where:
// ckID is a FOURCC that identifies the data contained in the chunk,
// ckSize is a 4-byte value representing the size of the data in ckData
//        (excluding the ckID FOURCC, the 4 bytes of ckSize, or the padding),
// ckData is 0 or more bytes of data, padded to the nearest WORD boundary
type chunk struct {
	ckID   FOURCC
	ckSize uint32
	ckData []byte
}

func (ck *chunk) String() string {
	return fmt.Sprintf("'%s' (%d)", string(ck.ckID[:]), ck.ckSize)
}

// A list has the form `'LIST' listSize listType listData, where:
// 'LIST' is the literal FOURCC 'LIST' (omitted from the struct),
// listSize is a 4-byte value representing the size of the contained data
//          (excluding the 'LIST' FOURCC and listSize),
// listType is a FOURCC code,
// listData is an array of chunks or lists, in any order
type list struct {
	listSize uint32
	listType FOURCC
	listData []riffElement
}

func (l *list) String() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "LIST '%4s' [%d] (\n", l.listType, l.listSize)
	for _, el := range l.listData {
		buf.WriteString("     ")
		buf.WriteString(el.String())
		buf.WriteString("\n")
	}
	buf.WriteString(")")
	return buf.String()
}

func DoStuff(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	var els []fmt.Stringer

	// iterators in go currently suck...
	s := Scan(file)
	for el, _ := s.Next(); el != nil; el, _ = s.Next() {
		els = append(els, el)
	}

	for _, el := range els {
		fmt.Println(el)
	}
}
