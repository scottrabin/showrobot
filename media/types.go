package media

import "time"

type MediaType int

type Movie struct {
	Name    string
	Runtime time.Duration
	Year    int
}

const MOVIE_DURATION = 60 * time.Minute

const (
	UNKNOWN MediaType = iota
	MOVIE
	TVSHOW
)
