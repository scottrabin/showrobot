package media

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/scottrabin/showrobot/media/riff"
	"os"
	"time"
)

var (
	ErrAviHeaderNotFound = errors.New("AVI Header not found in file")
)

type aviCodec struct{}

type aviMainHeader struct {
	//fcc                 FOURCC
	//cb                  uint32
	MicroSecPerFrame    uint32
	MaxBytesPerSec      uint32
	PaddingGranularity  uint32
	Flags               uint32
	TotalFrames         uint32
	InitialFrames       uint32
	Streams             uint32
	SuggestedBufferSize uint32
	Width               uint32
	Height              uint32
	Reserved            [4]uint32
}

// Decode fulfills the Codec interface
func (c *aviCodec) Decode(mf *MediaFile) (MediaInfo, error) {
	mi := MediaInfo{}
	file, err := os.Open(mf.Source)
	if err != nil {
		return mi, err
	}
	defer file.Close()

	// iterators in go currently suck...
	s := riff.Scan(file)
	avihChunk, err := c.findMainAviHeader(&s)
	if err != nil {
		return mi, err
	}

	var avih aviMainHeader
	buf := bytes.NewReader(avihChunk.Data)
	if err := binary.Read(buf, binary.LittleEndian, &avih); err != nil {
		return mi, err
	}

	mi.Duration = time.Duration(avih.MicroSecPerFrame) *
		time.Duration(avih.TotalFrames) * time.Microsecond

	return mi, nil
}

// findHeaderList finds the list with the ID "hdrl"
func (c *aviCodec) findHeaderList(s *riff.Scanner) (*riff.List, error) {
	for el, _ := s.Next(); el != nil; el, _ = s.Next() {
		if el, ok := el.(*riff.List); ok && el.Type.String() == "hdrl" {
			return el, nil
		}
	}

	return nil, ErrAviHeaderNotFound
}

// findMainHeader finds the riff chunk with the ID "avih"
func (c *aviCodec) findMainAviHeader(s *riff.Scanner) (*riff.Chunk, error) {
	hdrl, err := c.findHeaderList(s)
	if err != nil {
		return nil, err
	}

	for _, c := range hdrl.Data {
		if c, ok := c.(*riff.Chunk); ok && c.ID.String() == "avih" {
			return c, nil
		}
	}

	return nil, ErrAviHeaderNotFound
}

func init() {
	RegisterCodec(".avi", &aviCodec{})
}
