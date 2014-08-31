package media

import (
	"time"
)

type MediaType int

type Movie struct {
	Name    string
	Runtime time.Duration
	Year    int
}

type MediaInfo struct {
	Type     MediaType
	Duration time.Duration
}

const MOVIE_DURATION = 60 * time.Minute

const (
	UNKNOWN MediaType = iota
	MOVIE
	TVSHOW
)
