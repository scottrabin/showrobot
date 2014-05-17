package datasource

import "encoding/json"
import "github.com/scottrabin/showrobot/config"
import "github.com/scottrabin/showrobot/media"
import "io/ioutil"
import "net/http"
import "net/url"
import "time"

type TheMovieDB struct {
	config config.Configuration
}

func (ds *TheMovieDB) GetMovies(query string) []media.Movie {
	var jsonResults struct {
		Page    int
		Results []struct {
			Id           int
			Release_date string
			Title        string
		}
		Total_pages   int
		Total_results int
	}

	url := ds.getQuery(query)
	response, err := http.Get(url)
	if err != nil {
		return nil
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil
	}

	json.Unmarshal(body, &jsonResults)

	result := make([]media.Movie, len(jsonResults.Results))
	for i, r := range jsonResults.Results {
		released, _ := time.Parse("2006-01-02", r.Release_date)
		result[i] = media.Movie{Name: r.Title, Year: released.Year()}
	}

	return result
}

func (ds *TheMovieDB) getQuery(query string) string {
	apiKey, ok := ds.config.ApiKey["themoviedb"]
	if !ok {
		return ""
	}
	u := url.URL{}

	u.Scheme = "http"
	u.Host = "api.themoviedb.org"
	u.Path = "3/search/movie"

	qp := url.Values{}
	qp.Set("api_key", *apiKey)
	qp.Set("query", query)

	u.RawQuery = qp.Encode()

	return u.String()
}

func (ds *TheMovieDB) IsValid(mt media.MediaType) bool {
	switch mt {
	case media.MOVIE:
		_, ok := ds.config.ApiKey["themoviedb"]
		return ok
	default:
		return false
	}
}

func init() {
	datasources["themoviedb"] = func(conf config.Configuration) Datasource {
		return &TheMovieDB{conf}
	}
}
