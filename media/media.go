package media

import "fmt"
import "path/filepath"
import "regexp"
import "strings"
import "strconv"
import "time"

const MOVIE_DURATION = 60 * time.Minute

type MediaFile struct {
	source string
}

type Media interface {
	GetFileName() string
	GetRuntime() time.Duration
}

var year_locator_regexp, non_word_regexp *regexp.Regexp

func (m *MediaFile) GetFileName() string {
	return m.source
}

func init() {
	year_locator_regexp = regexp.MustCompile("(\\(\\d{4}\\)|\\[\\d{4}\\])")
	non_word_regexp = regexp.MustCompile("\\W+")
}

func NewMediaFile(file string) (m Media, err error) {
	ext := filepath.Ext(file)

	switch ext {
	case ".avi":
		m = &AviMediaFile{MediaFile{file}}
	default:
		err = fmt.Errorf("Unknown file extension: %s", ext)
	}

	return m, err
}

func getFileName(media Media) string {
	// start with the filename
	file := media.GetFileName()

	// remove the extension
	basename := filepath.Base(file)
	ext := filepath.Ext(file)
	return strings.TrimSuffix(basename, ext)
}

func GuessName(media Media) (guess string) {
	guess = getFileName(media)

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

func GuessYear(media Media) int {
	yearStr := year_locator_regexp.FindString(getFileName(media))
	if len(yearStr) > 0 {
		year, err := strconv.Atoi(yearStr)
		if err == nil {
			return year
		}
	}

	return 0
}

func IsMovie(media Media) bool {
	return media.GetRuntime() >= MOVIE_DURATION
}
