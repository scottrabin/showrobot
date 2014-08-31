package media

import (
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var year_locator_regexp, non_word_regexp *regexp.Regexp

func extractBaseName(mf *MediaFile) string {
	return strings.TrimSuffix(filepath.Base(mf.Source), filepath.Ext(mf.Source))
}

// GuessName attempts to guess the best searchable name for a given media file
// based on the filename
func GuessName(mf *MediaFile) string {
	guess := extractBaseName(mf)

	// Remove any year-looking parens
	// (include everything after that because it's usually junk metadata)
	guessBytes := []byte(guess)
	if loc := year_locator_regexp.FindIndex(guessBytes); loc != nil {
		guess = string(guessBytes[:loc[0]])
	}

	// split the string on non-word, non-number
	guess = strings.Join(non_word_regexp.Split(guess, -1), " ")

	return guess
}

// GuessYear attempts to guess the year associated with a media file based on
// the filename
func GuessYear(mf *MediaFile) int {
	yearStr := year_locator_regexp.FindString(extractBaseName(mf))
	if len(yearStr) > 0 {
		if year, err := strconv.Atoi(yearStr); err == nil {
			return year
		}
	}

	return 0
}

// GuessType attempts to guess the type of a media file (Movie or TV show)
// based on the duration of the media
func GuessType(mf *MediaFile) MediaType {
	info, err := mf.Info()
	switch {
	case err != nil:
		return UNKNOWN
	case info.Duration >= MOVIE_DURATION:
		return MOVIE
	default:
		return TVSHOW
	}
}

func init() {
	year_locator_regexp = regexp.MustCompile("(\\(\\d{4}\\)|\\[\\d{4}\\])")
	non_word_regexp = regexp.MustCompile("\\W+")
}
