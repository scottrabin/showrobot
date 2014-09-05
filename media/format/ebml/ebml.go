package ebml

import (
	"bytes"
	"encoding/binary"
	"errors"
	"time"
)

var (
	ErrInvalidFloatSize = errors.New("Invalid size for float; must be 4 or 8")
	ErrInvalidEBMLType  = errors.New("Invalid EBML type")
)

type ValueType int

type EBMLID uint64

type Element struct {
	// ID is the EBML ID of this element
	ID EBMLID
	// value is the raw slice of bytes contained in this element
	value []byte
	// parsedvalue is the cached return value of the Value() function
	parsedvalue interface{}
	// parseerror holds the error generated when attempting to parse the value
	parseerror error
}

type ElementMeta struct {
	Name string
	Type ValueType
}

const (
	typeGuess ValueType = iota
	typeMaster
	typeUnsignedInt
	typeSignedInt
	typeString
	typeUtf8String
	typeBinary
	typeFloat
	typeDate
)

// Value parses the value contained in the raw byte slice of the element into
// the proper value; this value is cached so successive calls to Value() do not
// reparse the byte slice each time.
func (el *Element) Value() (interface{}, error) {
	if el.parsedvalue == nil && el.parseerror == nil {
		switch el.Meta().Type {
		case typeMaster:
			el.parsedvalue, _ = ReadElements(el.value)
		case typeUnsignedInt:
			el.parsedvalue, _ = binary.Uvarint(el.value)
		case typeSignedInt:
			el.parsedvalue, _ = binary.Varint(el.value)
		case typeString, typeUtf8String:
			el.parsedvalue = string(el.value)
		case typeBinary, typeGuess:
			el.parsedvalue = el.value
		case typeFloat:
			switch len(el.value) {
			case 4:
				var v float32
				el.parseerror = binary.Read(bytes.NewBuffer(el.value), binary.BigEndian, &v)
				el.parsedvalue = v
			case 8:
				var v float64
				el.parseerror = binary.Read(bytes.NewBuffer(el.value), binary.BigEndian, &v)
				el.parsedvalue = float64(0)
			default:
				el.parseerror = ErrInvalidFloatSize
			}
		case typeDate:
			// Date - signed 8 octets integer in nanoseconds with 0 indicating the
			// precise beginning of the millennium (at 2001-01-01T00:00:00,000000000 UTC)
			var (
				rv        int64
				startDate = time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)
			)
			el.parseerror = binary.Read(bytes.NewBuffer(el.value), binary.BigEndian, &rv)
			el.parsedvalue = startDate.Add(time.Duration(rv) * time.Nanosecond)
		default:
			el.parseerror = ErrInvalidEBMLType
		}
	}

	return el.parsedvalue, el.parseerror
}

// Meta returns the element type metadata associated with the element
func (el *Element) Meta() ElementMeta {
	if meta, known := ebmlIdMap[el.ID]; known {
		return meta
	}

	return ebmlIdMap[ElementUnknown]
}

// ReadID takes a byte slice and reads the next ID from the beginning.
// Returns the next ID and the length of the ID, in bytes
func ReadId(s []byte) (EBMLID, int) {
	for i, b := range []byte{0x80, 0x40, 0x20, 0x10, 0x08, 0x04, 0x02, 0x01} {
		if s[0]&b > 0 {
			var rv EBMLID
			for _, p := range s[0 : i+1] {
				rv = (rv << 8) + EBMLID(p)
			}
			return rv, i + 1
		}
	}

	return 0, 0
}

// ReadSize takes a byte slice of an EBML document and reads the next size
// parameter from the beginning. Returns the size and the length of the size, in bytes
func ReadSize(s []byte) (uint64, int) {
	// the first non-zero bit in the byte slice (interpreted as big-endian)
	// indicates the length of the integer; the number of leading zeroes + 1
	// is the length, in octets, of the number
	for i, b := range []byte{0x80, 0x40, 0x20, 0x10, 0x08, 0x04, 0x02, 0x01} {
		if s[0]&b > 0 {
			var rv uint64 = uint64(s[0] ^ b)
			if i > 0 {
				for _, p := range s[1 : i+1] {
					rv = (rv << 8) + uint64(p)
				}
			}
			return rv, i + 1
		}
	}

	return 0, 0
}

// ReadElement takes a byte slice of an EBML document and reads the next element,
// including its ID, size, and value, from the beginning of the slice. Returns
// the element and the length of the element in bytes
func ReadElement(s []byte) (Element, int64) {
	var (
		ret          Element
		idlen, szlen int
		sz           uint64
	)

	if len(s) > 0 {
		ret.ID, idlen = ReadId(s)
		sz, szlen = ReadSize(s[idlen:])
		ret.value = s[idlen+szlen : uint64(idlen+szlen)+sz]
	}

	return ret, (int64(idlen+szlen) + int64(sz))
}

// ReadElements takes a byte slice and reads all of the elements contained in it
func ReadElements(s []byte) ([]Element, int64) {
	var (
		els    []Element
		offset int64
	)

	for length := int64(len(s)); offset < length; {
		el, size := ReadElement(s[offset:])
		els = append(els, el)
		offset += size
	}

	return els, offset
}
