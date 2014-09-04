package media

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"strings"
	"syscall"
	"time"
)

type ebmlType int

const (
	typeGuess ebmlType = iota
	typeMaster
	typeUnsignedInt
	typeSignedInt
	typeString
	typeUtf8String
	typeBinary
	typeFloat
	typeDate
)

type ebmlElement struct {
	ID    uint64
	Type  ebmlType
	Value interface{}
}

type ebmlElementMeta struct {
	Name string
	Type ebmlType
}

var ebmlIdMap = map[uint64]ebmlElementMeta{
	0: ebmlElementMeta{
		Name: "Unknown", // not part of spec
		Type: typeGuess,
	},

	0xEC: ebmlElementMeta{
		Name: "Void",
		Type: typeBinary,
	},
	0xBF: ebmlElementMeta{
		Name: "CRC-32",
		Type: typeBinary,
	},

	0x1A45DFA3: ebmlElementMeta{
		Name: "EBML",
		Type: typeMaster,
	},
	0x4286: ebmlElementMeta{
		Name: "EBMLVersion",
		Type: typeUnsignedInt,
	},
	0x42F7: ebmlElementMeta{
		Name: "EBMLReadVersion",
		Type: typeUnsignedInt,
	},
	0x42F2: ebmlElementMeta{
		Name: "EBMLMaxIDLength",
		Type: typeUnsignedInt,
	},
	0x42F3: ebmlElementMeta{
		Name: "EBMLMaxSizeLength",
		Type: typeUnsignedInt,
	},
	0x4282: ebmlElementMeta{
		Name: "DocType",
		Type: typeString,
	},
	0x4287: ebmlElementMeta{
		Name: "DocTypeVersion",
		Type: typeUnsignedInt,
	},
	0x4285: ebmlElementMeta{
		Name: "DocTypeReadVersion",
		Type: typeUnsignedInt,
	},

	0x18538067: ebmlElementMeta{
		Name: "Segment",
		Type: typeMaster,
	},

	0x114D9B74: ebmlElementMeta{
		Name: "SeekHead",
		Type: typeMaster,
	},
	0x4DBB: ebmlElementMeta{
		Name: "Seek",
		Type: typeMaster,
	},
	0x53AB: ebmlElementMeta{
		Name: "SeekID",
		Type: typeBinary,
	},
	0x53AC: ebmlElementMeta{
		Name: "SeekPosition",
		Type: typeUnsignedInt,
	},

	0x1549A966: ebmlElementMeta{
		Name: "Info",
		Type: typeMaster,
	},
	0x73A4: ebmlElementMeta{
		Name: "SegmentUID",
		Type: typeBinary,
	},
	0x7384: ebmlElementMeta{
		Name: "SegmentFilename",
		Type: typeString, // UTF-8
	},
	0x3CB923: ebmlElementMeta{
		Name: "PrevUID",
		Type: typeBinary,
	},
	0x3C83AB: ebmlElementMeta{
		Name: "PrevFilename",
		Type: typeString, // UTF-8
	},
	0x3EB923: ebmlElementMeta{
		Name: "NextUID",
		Type: typeBinary,
	},
	0x3E83BB: ebmlElementMeta{
		Name: "NextFilename",
		Type: typeString, // UTF-8
	},
	0x4444: ebmlElementMeta{
		Name: "SegmentFamily",
		Type: typeBinary,
	},
	0x6924: ebmlElementMeta{
		Name: "ChapterTranslate",
		Type: typeMaster,
	},
	0x69FC: ebmlElementMeta{
		Name: "ChapterTranslateEditionUID",
		Type: typeUnsignedInt,
	},
	0x69BF: ebmlElementMeta{
		Name: "ChapterTranslateCodec",
		Type: typeUnsignedInt,
	},
	0x69A5: ebmlElementMeta{
		Name: "ChapterTranslateID",
		Type: typeBinary,
	},
	0x2AD7B1: ebmlElementMeta{
		Name: "TimecodeScale",
		Type: typeUnsignedInt,
	},
	0x4489: ebmlElementMeta{
		Name: "Duration",
		Type: typeFloat,
	},
	0x4461: ebmlElementMeta{
		Name: "DateUTC",
		Type: typeDate,
	},
	0x7BA9: ebmlElementMeta{
		Name: "Title",
		Type: typeString, // UTF-8
	},
	0x4D80: ebmlElementMeta{
		Name: "MuxingApp",
		Type: typeString, // UTF-8
	},
	0x5741: ebmlElementMeta{
		Name: "WritingApp",
		Type: typeString, // UTF-8
	},

	0x1F43B675: ebmlElementMeta{
		Name: "Cluster",
		Type: typeMaster,
	},
	0xE7: ebmlElementMeta{
		Name: "Timecode",
		Type: typeUnsignedInt,
	},
	0x5854: ebmlElementMeta{
		Name: "SilentTracks",
		Type: typeMaster,
	},
	0x58D7: ebmlElementMeta{
		Name: "SilentTrackNumber",
		Type: typeUnsignedInt,
	},
	0xA7: ebmlElementMeta{
		Name: "Position",
		Type: typeUnsignedInt,
	},
	0xAB: ebmlElementMeta{
		Name: "PrevSize",
		Type: typeUnsignedInt,
	},
	0xA3: ebmlElementMeta{
		Name: "SimpleBlock",
		Type: typeBinary,
	},
	0xA0: ebmlElementMeta{
		Name: "BlockGroup",
		Type: typeMaster,
	},
	0xA1: ebmlElementMeta{
		Name: "Block",
		Type: typeBinary,
	},
}

func (el *ebmlElement) Meta() ebmlElementMeta {
	if meta, known := ebmlIdMap[el.ID]; known {
		return meta
	}
	return ebmlIdMap[0]
}

type matroskaCodec struct{}

func (m *matroskaCodec) Decode(mf *MediaFile) (MediaInfo, error) {
	mi := MediaInfo{}

	file, err := os.Open(mf.Source)
	if err != nil {
		return mi, err
	}
	fileinfo, err := file.Stat()
	if err != nil {
		return mi, err
	}

	mmap, err := syscall.Mmap(int(file.Fd()), 0, int(fileinfo.Size()),
		syscall.PROT_READ, syscall.MAP_SHARED)
	if err != nil {
		return mi, err
	}

	els, _ := m.process(mmap)
	m.PrintElements(els, 0)

	return mi, nil
}

func (m *matroskaCodec) PrintElements(els []ebmlElement, depth int) {
	for _, el := range els {
		switch v := el.Value.(type) {
		case []ebmlElement:
			fmt.Printf("%s%x (%s)\n", strings.Repeat(" ", depth*4), el.ID, el.Meta().Name)
			if el.ID != 0x1F43B675 {
				// skip clusters; for everything else, show children
				m.PrintElements(v, depth+1)
			}
		case []byte:
			if len(v) > 24 {
				v = v[:24]
			}
			fmt.Printf("%s%x (%s)\t%v\n", strings.Repeat(" ", depth*4), el.ID, el.Meta().Name, v)
		default:
			fmt.Printf("%s%x (%s)\t%v\n",
				strings.Repeat(" ", depth*4), el.ID, el.Meta().Name, v)
		}
	}
}

func (m *matroskaCodec) readId(s []byte) (uint64, int) {
	for i, b := range []byte{0x80, 0x40, 0x20, 0x10, 0x08, 0x04, 0x02, 0x01} {
		if s[0]&b > 0 {
			var rv uint64
			for _, p := range s[0 : i+1] {
				rv = (rv << 8) + uint64(p)
			}
			return rv, i + 1
		}
	}

	return 0, 0
}

func (m *matroskaCodec) readSize(s []byte) (uint64, int) {
	// the first non-zero bit in the byte slice (interpreted as big-endian)
	// indicates the length of the integer; the number of leading zeroes + 1
	// is the length, in octets, of the number
	for i, b := range []byte{0x80, 0x40, 0x20, 0x10, 0x08, 0x04, 0x02, 0x01} {
		if s[0]&b > 0 {
			buf := make([]byte, i+1)
			copy(buf, s[0:i+1])
			buf[0] = buf[0] ^ b

			var rv uint64
			for _, p := range buf {
				rv = (rv << 8) + uint64(p)
			}
			return rv, i + 1
		}
	}

	return 0, 0
}

func (m *matroskaCodec) process(s []byte) ([]ebmlElement, int64) {
	var off int64
	els := make([]ebmlElement, 0)
	for {
		el, n := m.readElement(s[off:])
		if n == 0 {
			break
		}
		els = append(els, el)
		off = off + n

	}

	return els, off
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (m *matroskaCodec) readElement(s []byte) (ebmlElement, int64) {
	var (
		ret          ebmlElement
		idlen, szlen int
		sz           uint64
		//err          error
	)

	if len(s) > 0 {
		ret.ID, idlen = m.readId(s)
		sz, szlen = m.readSize(s[idlen:])

		ret.Value, _ = m.readValue(s[idlen+szlen:uint64(idlen+szlen)+sz], ret.Meta().Type)
	}
	return ret, (int64(idlen+szlen) + int64(sz))
}

func (m *matroskaCodec) readValue(s []byte, typ ebmlType) (interface{}, error) {
	switch typ {
	case typeMaster:
		el, n := m.process(s)
		if n == 0 {
			return nil, fmt.Errorf("Could not process value")
		}
		return el, nil
	case typeUnsignedInt:
		v, _ := binary.Uvarint(s)
		return v, nil
	case typeSignedInt:
		v, _ := binary.Varint(s)
		return v, nil
	case typeString, typeUtf8String:
		return string(s), nil
	case typeBinary:
		return s, nil
	case typeFloat:
		if len(s) == 4 {
			var rv float32
			err := binary.Read(bytes.NewBuffer(s), binary.BigEndian, &rv)
			return rv, err
		} else if len(s) == 8 {
			var rv float64
			err := binary.Read(bytes.NewBuffer(s), binary.BigEndian, &rv)
			return rv, err
		}
		return float32(0), fmt.Errorf("Cannot convert odd-sized value %x to float", s)
	case typeDate:
		// Date - signed 8 octets integer in nanoseconds with 0 indicating the
		// precise beginning of the millennium (at 2001-01-01T00:00:00,000000000 UTC)
		var (
			rv        int64
			startDate = time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)
		)
		err := binary.Read(bytes.NewBuffer(s), binary.BigEndian, &rv)
		return startDate.Add(time.Duration(rv) * time.Nanosecond), err

		//v, _ := binary.Varint(s)
		//return time.Unix(0, v), nil
	case typeGuess:
		if len(s) > 8 {
			return s, nil
		}

		v, _ := binary.Uvarint(s)
		return v, nil
	}

	return nil, fmt.Errorf("Invalid ebml type: %d", typ)
}

func init() {
	RegisterCodec(".mkv", &matroskaCodec{})
}
