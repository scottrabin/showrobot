package datasource

import "github.com/scottrabin/showrobot/config"
import "github.com/scottrabin/showrobot/media"

type MovieDatasource interface {
	GetMovies(media.Media) []media.Movie
}

func NewMovieSource(conf config.Configuration, which string) MovieDatasource {
	switch which {
	case "themoviedb":
		return &TheMovieDB{conf}
	}

	return nil
}
