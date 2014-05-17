package media

import "fmt"
import "path/filepath"
import "time"

// A Media is a file representing a specific piece of content (Movie or TV show)
type Media interface {
	Source() string
	Format() MediaFormat
	Duration() time.Duration
}

// A MediaFormat is an adapter for determining metadata about a Media that
// cannot be deduced from the filename alone, such as the duration of the content
type MediaFormat interface {
	Duration(Media) time.Duration
}

// A MediaFile is a file on disk that contains a piece of content (Movie or
// TV Show)
type MediaFile struct {
	source string
	format MediaFormat
}

func (m *MediaFile) Source() string {
	return m.source
}

func (m *MediaFile) Format() MediaFormat {
	return m.format
}

func (m *MediaFile) Duration() time.Duration {
	// this should be overridden by all structs that embed the MediaFile type
	return m.Format().Duration(m)
}

func New(file string) (Media, error) {
	var err error

	unk, _ := mediaFormats[""]
	mf := &MediaFile{file, unk}
	ext := filepath.Ext(file)

	format, supported := mediaFormats[ext]

	if supported {
		mf.format = format
	} else {
		err = fmt.Errorf("unknown file extension: %s", ext)
	}

	return mf, err
}

var mediaFormats = make(map[string]MediaFormat)

func Register(ext string, f MediaFormat) {
	mediaFormats[ext] = f
}
