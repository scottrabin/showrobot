package media

import "fmt"
import "path/filepath"
import "regexp"
import "strings"
import "strconv"
import "time"

type Movie struct {
	Name    string
	Runtime time.Duration
	Year    int
}

type MediaType int

type MediaFile struct {
	source string
}

type Media interface {
	GetPath() string
	GetFileName() string
	GetExtension() string
	GuessName() string
	GuessYear() int
	GuessType() MediaType
	GetRuntime() time.Duration
}

const MOVIE_DURATION = 60 * time.Minute

const (
	UNKNOWN MediaType = iota
	MOVIE
	TVSHOW
)

var year_locator_regexp, non_word_regexp *regexp.Regexp

func NewMedia(file string) (m Media, err error) {
	mf := MediaFile{file}

	switch mf.GetExtension() {
	case ".avi":
		m = &AviMediaFile{mf}
	default:
		err = fmt.Errorf("Unknown file extension: %s", mf.GetExtension())
	}

	return m, err
}

func (m *MediaFile) GetPath() string {
	return m.source
}

func (m *MediaFile) GetExtension() string {
	return filepath.Ext(m.GetPath())
}

func (m *MediaFile) GetFileName() string {
	// remove the extension
	basename := filepath.Base(m.GetPath())
	return strings.TrimSuffix(basename, m.GetExtension())
}

func (m *MediaFile) GuessName() (guess string) {
	guess = m.GetFileName()

	// Remove any year-looking parens
	// (include everything after that because it's usually junk metadata)
	guessBytes := []byte(guess)
	if loc := year_locator_regexp.FindIndex(guessBytes); loc != nil {
		guess = string(guessBytes[:loc[0]])
	}

	// split the string on non-word, non-number
	guess = strings.Join(non_word_regexp.Split(guess, -1), " ")

	return
}

func (m *MediaFile) GuessYear() int {
	yearStr := year_locator_regexp.FindString(m.GetFileName())
	if len(yearStr) > 0 {
		year, err := strconv.Atoi(yearStr)
		if err == nil {
			return year
		}
	}

	return 0
}

func (m *MediaFile) GetRuntime() time.Duration {
	// this should be overridden by all structs that embed the MediaFile type
	return time.Duration(0)
}

func (m *MediaFile) GuessType() MediaType {
	if m.GetRuntime() >= MOVIE_DURATION {
		return MOVIE
	} else {
		return TVSHOW
	}
}

func init() {
	year_locator_regexp = regexp.MustCompile("(\\(\\d{4}\\)|\\[\\d{4}\\])")
	non_word_regexp = regexp.MustCompile("\\W+")
}
