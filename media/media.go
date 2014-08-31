package media

import (
	"path/filepath"
)

// A MediaFile is a file on disk that contains a piece of content (Movie or
// TV Show)
type MediaFile struct {
	Source string
	info   *MediaInfo
}

// NewFile creates a new MediaFile object that represents the given file
func NewFile(file string) *MediaFile {
	return &MediaFile{
		Source: file,
	}
}

// Parse uses the given codec to decode the media file and read its information
// If no codec is specified, the codec registered to the file's extension
// will be used instead.
func (mf *MediaFile) Parse(c *Codec) error {
	if c == nil {
		var supported bool
		if c, supported = codecs[filepath.Ext(mf.Source)]; !supported {
			return ErrNoCodec
		}
	}

	if info, err := c.Decode(mf); err == nil {
		mf.info = &info
	} else {
		return err
	}

	return nil
}

// Info retrieves the media info associated with the media file
// If the media file has not been decoded, Info will first call Parse(nil)
// on the media file
func (mf *MediaFile) Info() (*MediaInfo, error) {
	var err error

	if mf.info == nil {
		err = mf.Parse(nil)
	}

	return mf.info, err
}
