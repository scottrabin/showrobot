package datasource

import "fmt"
import "github.com/scottrabin/showrobot/config"
import "github.com/scottrabin/showrobot/media"

type Datasource interface {
	GetMovies(string) []media.Movie
	IsValid(media.MediaType) bool
}
type DatasourceCreator func(config.Configuration) Datasource

var datasources = make(map[string]DatasourceCreator)

func New(conf config.Configuration, mt media.MediaType) (Datasource, error) {
	for _, mkSrc := range datasources {
		source := mkSrc(conf)
		if source.IsValid(mt) {
			return source, nil
		}
	}
	return nil, fmt.Errorf("no valid sources configured")
}
