package datasource

import "github.com/scottrabin/showrobot/config"
import "github.com/scottrabin/showrobot/media"
import "time"

type Movie struct {
	Name    string
	Runtime time.Duration
	Year    int
}

type MovieDatasource interface {
	GetMovies(media.Media) []Movie
}

func NewMovieSource(conf config.Configuration, which string) MovieDatasource {
	switch which {
	case "themoviedb":
		return &TheMovieDB{conf}
	}

	return nil
}
