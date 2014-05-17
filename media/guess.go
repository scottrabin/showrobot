package media

import "path/filepath"
import "regexp"
import "strconv"
import "strings"

var year_locator_regexp, non_word_regexp *regexp.Regexp

func extractBaseName(m Media) string {
	src := m.Source()
	basename := filepath.Base(src)
	extension := filepath.Ext(src)

	return strings.TrimSuffix(basename, extension)
}

// GuessName attempts to guess the best searchable name for a given media file
// based on the filename
func GuessName(m Media) string {
	guess := extractBaseName(m)

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
func GuessYear(m Media) int {
	fn := extractBaseName(m)
	yearStr := year_locator_regexp.FindString(fn)
	if len(yearStr) > 0 {
		year, err := strconv.Atoi(yearStr)
		if err == nil {
			return year
		}
	}

	return 0
}

// GuessType attempts to guess the type of a media file (Movie or TV show)
// based on the duration of the media
func GuessType(m Media) MediaType {
	if m.Duration() >= MOVIE_DURATION {
		return MOVIE
	} else {
		return TVSHOW
	}
}

func init() {
	year_locator_regexp = regexp.MustCompile("(\\(\\d{4}\\)|\\[\\d{4}\\])")
	non_word_regexp = regexp.MustCompile("\\W+")
}
