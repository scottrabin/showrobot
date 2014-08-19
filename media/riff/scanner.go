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
func (s *Scanner) Next() (riffElement, error) {
	el, _, err := s.next()
	return el, err
}

// next reads the next element and returns the amount read from the wrapped
// io.Reader
func (s *Scanner) next() (riffElement, int, error) {
	var (
		fcc        FOURCC
		el         riffElement
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
func (s *Scanner) scanHeader() (*riffheader, int, error) {
	var (
		hdr          riffheader
		fsLen, ftLen int
		err          error
	)

	if hdr.fileSize, fsLen, err = s.scanSize(); err != nil {
		return nil, fsLen, err
	}

	if hdr.fileType, ftLen, err = s.scanFcc(); err != nil {
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
func (s *Scanner) scanList() (*list, int, error) {
	var (
		l                   list
		lsLen, ltLen, ldLen int
		err                 error
	)

	if l.listSize, lsLen, err = s.scanSize(); err != nil {
		return nil, lsLen, err
	}

	if l.listType, ltLen, err = s.scanFcc(); err != nil {
		return nil, (lsLen + ltLen), err
	}

	// the listSize includes the listType length, and listSize
	// should be equal to the length of the listType and listData
	for ldLen < int(l.listSize)-ltLen {
		el, read, err := s.next()
		if err != nil {
			return nil, (lsLen + ltLen + ldLen), err
		}

		l.listData = append(l.listData, el)
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
func (s *Scanner) scanChunk(id FOURCC) (*chunk, int, error) {
	var (
		c              chunk
		csLen, dataLen int
		err            error
	)
	c.ckID = id

	if c.ckSize, csLen, err = s.scanSize(); err != nil {
		return nil, csLen, err
	}

	if c.ckData, dataLen, err = s.scanData(c.ckSize); err != nil {
		return nil, csLen + dataLen, err
	}

	// chunks are padded to the nearest WORD boundary, so skip that many bytes
	pad := c.ckSize % 4
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
