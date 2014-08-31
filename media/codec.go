package media

import (
	"errors"
)

type Codec struct {
	Decode func(*MediaFile) (MediaInfo, error)
}

var codecs = make(map[string]*Codec)
var ErrNoCodec = errors.New("An appropriate codec could not be found")

// Register attaches a file extension (e.g. ".avi") to a specific
// codec that can retrieve MediaInfo from a given MediaFile
func RegisterCodec(ext string, c *Codec) {
	codecs[ext] = c
}

// LookupCodec retrieves a codec registered for a the given file extension
// if no codec is registered for the given extension, returns ErrNoCodec
func LookupCodec(ext string) (*Codec, error) {
	if codec, ok := codecs[ext]; ok {
		return codec, nil
	}

	return nil, ErrNoCodec
}
