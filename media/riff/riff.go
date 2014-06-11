package riff

import "encoding/binary"
import "fmt"
import "io"
import "io/ioutil"
import "math"
import "os"
import "strings"

type FOURCC [4]byte

type riffheader struct {
	fileSize [4]byte
	fileType FOURCC
	data []byte
}

type chunk struct {
	ckID FOURCC
	ckSize [4]byte
	ckData []byte
}

type list struct {
	listSize [4]byte
	listType FOURCC
	listData []interface{} // TODO this can be lists or chunks
}

const DATA_READ_LIMIT = 256
// TODO Figure out a way to make this a const
var RIFF_RIFF = FOURCC{'R', 'I', 'F', 'F'}
var RIFF_LIST = FOURCC{'L', 'I', 'S', 'T'}

func (h *riffheader) String() string {
	return fmt.Sprintf("RIFF '%s' (%s)", string(h.fileType[:]), string(h.data))
}

func (ck *chunk) String() string {
	if len(ck.ckData) > 0 {
		return fmt.Sprintf("'%s' [%s]", string(ck.ckID[:]), string(ck.ckData))
	}
	return fmt.Sprintf("'%s' (%d)", string(ck.ckID[:]), convertSize(ck.ckSize))
}

func (l *list) String() string {
	return fmt.Sprintf("LIST (\n%s)\n", l.listData...)
}

func convertSize(s [4]byte) int64 {
	r, _ := binary.Uvarint(s[:])
	return int64(r)
}

// Pad the size of the data for various RIFF elements to the
// nearest WORD boundary
func padSize(sz int64) int64 {
	return 4 * int64(math.Floor(float64(sz) / 4.0))
}

// Read the "data" field from a RIFF element if it isn't too big
func readData(r io.Reader, n int64) []byte {
	var dst []byte
	if n <= DATA_READ_LIMIT {
		dst := make([]byte, n, n)
		r.Read(dst)
	} else {
		io.CopyN(ioutil.Discard, r, n)
	}
	return dst
}

func readNext(r io.Reader) (interface{}, error) {
	// TODO the first return value can be a `list` or a `chunk`
	var fcc FOURCC

	n, err := r.Read(fcc[:])
	if n == 0 || err != nil {
		return nil, err
	}

	switch fcc {
	case RIFF_RIFF:
		h := riffheader{}
		// read in the file size and type
		r.Read(h.fileSize[:])
		r.Read(h.fileType[:])
		//h.data = readData(r, convertSize(h.fileSize))
		return h, nil
	case RIFF_LIST:
		l := list{}
		r.Read(l.listSize[:])
		r.Read(l.listType[:])
		return l, nil
	default:
		ck := chunk{}
		ck.ckID = fcc
		r.Read(ck.ckSize[:])
		sz := convertSize(ck.ckSize)
		ck.ckData = readData(r, sz)
		// chunks are padded to the nearest WORD boundary, so skip that many bytes
		io.CopyN(ioutil.Discard, r, padSize(sz) - sz)
		return ck, nil
	}
}

func DoStuff(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	doRead(file, 0)
}

func doRead(r io.Reader, indent int) {
	for {
		p, err := readNext(r)
		if err != nil {
			if err != io.EOF {
				fmt.Println("ERROR:", err)
			}
			break
		}

		switch p := p.(type) {
		case riffheader:
			fmt.Println(p.String())
		case list:
			dataLen := convertSize(p.listSize)
			if dataLen < DATA_READ_LIMIT {
				fmt.Printf("%sLIST '%4s' [%d] (\n", strings.Repeat("    ", indent), string(p.listType[:]), dataLen)
				doRead(r, indent + 1)
				fmt.Printf("%s)\n", strings.Repeat("    ", indent))
			} else {
				fmt.Printf("%sLIST '%4s' [%d, truncated]\n", strings.Repeat("    ", indent), string(p.listType[:]), dataLen)
				io.CopyN(ioutil.Discard, r, dataLen - 4)
			}

		case chunk:
			fmt.Printf("%s%s\n", strings.Repeat("    ", indent), p.String())
		}
	}
}
