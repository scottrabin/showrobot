package riff

import (
	"encoding/binary"
	"io"
	"io/ioutil"
)

// Maximum size of data to read into an element
const DATA_READ_LIMIT = 256

// A Scanner wraps an io.Reader and yields riff elements sequentially as
// needed to fulfill successive calls to Next()
type Scanner struct {
	r io.Reader
}

// Scan wraps an io.Reader to allow reading successive RIFF elements
func Scan(r io.Reader) Scanner {
	return Scanner{r: r}
}

// Next attempts to read the next RIFF element from the wrapped io.Reader
func (s *Scanner) Next() (RiffElement, error) {
	el, _, err := s.next()
	return el, err
}

// next reads the next element and returns the amount read from the wrapped
// io.Reader
func (s *Scanner) next() (RiffElement, int, error) {
	var (
		fcc        FOURCC
		el         RiffElement
		flen, elen int
		err        error
	)

	if fcc, flen, err = s.scanFcc(); err != nil {
		return nil, flen, err
	}

	switch string(fcc[:]) {
	case "RIFF":
		el, elen, err = s.scanHeader()
	case "LIST":
		el, elen, err = s.scanList()
	default:
		el, elen, err = s.scanChunk(fcc)
	}

	return el, (flen + elen), err
}

// scanHeader reads the header element from the beginning of the io.Reader,
// not including the initial `RIFF` FOURCC that indicates the beginning of
// a RIFF header. The RIFF header has the following form:
// 'RIFF' fileSize fileType (data)
// where:
//     'RIFF' is the literal FOURCC code 'RIFF'
//     fileSize is a 4-byte unsigned integer indicating the size of the file
//     filetype is a FOURCC that identifies the specific file type
func (s *Scanner) scanHeader() (*RiffHeader, int, error) {
	var (
		hdr          RiffHeader
		fsLen, ftLen int
		err          error
	)

	if hdr.Size, fsLen, err = s.scanSize(); err != nil {
		return nil, fsLen, err
	}

	if hdr.Type, ftLen, err = s.scanFcc(); err != nil {
		return nil, (fsLen + ftLen), err
	}

	return &hdr, (fsLen + ftLen), nil
}

// scanList reads the next list element from the wrapped io.Reader, not
// including the `LIST` FOURCC that indicates the start of the list. A
// RIFF LIST element has the following form:
// 'LIST' listSize listType listData
// where:
//     'LIST' is the literal FOURCC code 'LIST'
//     listSize is a 4-byte unsigned integer indicating the size of the list
//     listType is a FOURCC code
//     listData is a series of lists or chunks, in any order
func (s *Scanner) scanList() (*List, int, error) {
	var (
		l                   List
		lsLen, ltLen, ldLen int
		err                 error
	)

	if l.Size, lsLen, err = s.scanSize(); err != nil {
		return nil, lsLen, err
	}

	if l.Type, ltLen, err = s.scanFcc(); err != nil {
		return nil, (lsLen + ltLen), err
	}

	// the listSize includes the listType length, and listSize
	// should be equal to the length of the listType and listData
	for ldLen < int(l.Size)-ltLen {
		el, read, err := s.next()
		if err != nil {
			return nil, (lsLen + ltLen + ldLen), err
		}

		l.Data = append(l.Data, el)
		ldLen = ldLen + read
	}

	return &l, (lsLen + ltLen + ldLen), nil
}

// scanChunk reads the next chunk from the wrapped io.Reader, not including
// the initial FOURCC code that is neither `RIFF` nor `LIST`, and should be
// passed in to the scanChunk function as the `id` parameter. A chunk has
// the following form:
// ckID ckSize ckData
// where:
//     ckID is a FOURCC that identifies the data contained in the chunk
//     ckSize is a 4-byte unsigned integer indicating the size of the data
//     ckData is zero or more bytes of data, padded to the nearest WORD boundary
func (s *Scanner) scanChunk(id FOURCC) (*Chunk, int, error) {
	var (
		c              Chunk
		csLen, dataLen int
		err            error
	)
	c.ID = id

	if c.Size, csLen, err = s.scanSize(); err != nil {
		return nil, csLen, err
	}

	if c.Data, dataLen, err = s.scanData(c.Size); err != nil {
		return nil, csLen + dataLen, err
	}

	// chunks are padded to the nearest WORD boundary, so skip that many bytes
	pad := c.Size % 4
	if pad > 0 {
		io.CopyN(ioutil.Discard, s.r, int64(4-pad))
	}
	return &c, (csLen + dataLen + int(pad)), err
}

// scanFcc scans the wrapped io.Reader for the next FOURCC element in the stream.
// Returns the FOURCC and the number of bytes read.
func (s *Scanner) scanFcc() (FOURCC, int, error) {
	var fcc FOURCC
	n, err := s.r.Read(fcc[:])
	return fcc, n, err
}

// scanSize scans the wrapped io.Reader for the next four bytes as a uint32
// Returns the scanned number and the number of bytes read
func (s *Scanner) scanSize() (uint32, int, error) {
	var rv uint32
	err := binary.Read(s.r, binary.LittleEndian, &rv)
	return rv, 4, err
}

// scanData reads the "data" field for a RIFF chunk, if it isn't too big
func (s *Scanner) scanData(n uint32) ([]byte, int, error) {
	if n > DATA_READ_LIMIT {
		read64, err := io.CopyN(ioutil.Discard, s.r, int64(n))
		return make([]byte, 0), int(read64), err
	}

	dst := make([]byte, n)
	read, err := s.r.Read(dst)

	return dst, read, err
}
