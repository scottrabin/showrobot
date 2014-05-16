package main

import "bytes"
import "fmt"
import "github.com/codegangsta/cli"
import "github.com/scottrabin/showrobot/config"
import "github.com/scottrabin/showrobot/datasource"
import "github.com/scottrabin/showrobot/media"
import "os"
import "text/template"

func main() {
	fileFlag := cli.StringFlag{"file, f", config.GetDefaultConfigurationPath(), "Use the specified configuration file"}
	typeFlag := cli.StringFlag{"type, t", "auto", "Use the specified type for matching the media file"}
	queryFlag := cli.StringFlag{"query, q", "", "Use the specified string as the search term instead of guessing from the file name"}

	app := cli.NewApp()
	app.Name = "showrobot"
	app.Version = "0.0.1"
	app.Usage = "Classify & rename media files using remote datasources"
	app.Commands = []cli.Command{
		{
			Name:  "config",
			Usage: "get or set the configuration value for a given key",
			Description: `Modify the default configuration values for showrobot
when run without command line overrides`,
			Flags: []cli.Flag{fileFlag},
			Action: func(c *cli.Context) {
				var err error
				var value string
				args := c.Args()
				conf := loadConfig(c)

				switch len(args) {
				case 2:
					// set the configuration value
					err = conf.Set(args[0], args[1])
					if err == nil {
						err = conf.Save(c.String("file"))
					}
				case 1:
					// get the configuration value
					value, err = conf.Get(args[0])
					if err == nil {
						fmt.Println(value)
					}
				default:
					err = fmt.Errorf(
						"showrobot config requires 1 or 2 arguments; got %d", len(args))
				}

				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
			},
		},
		{
			Name:        "identify",
			Usage:       "display best media match for given file",
			Description: "Report the best matching media item for the given file",
			Flags:       []cli.Flag{fileFlag, typeFlag, queryFlag},
			Action: func(c *cli.Context) {
				args := c.Args()
				conf := loadConfig(c)

				m, err := media.New(args[0])
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				var matches []media.Movie

				switch mtype, err := getMediaType(c, m); mtype {
				case media.MOVIE:
					ds := datasource.NewMovieSource(conf, "themoviedb")
					query := getQuery(c, m)
					matches = ds.GetMovies(query)
				case media.TVSHOW:
					// TODO
					err = fmt.Errorf("TV show identification not yet implemented")
					fallthrough
				case media.UNKNOWN:
					fmt.Println(err)
					os.Exit(1)
				}

				// TODO sort movies by match likelihood

				// TODO handle errors here
				tmpl, _ := template.New("movie").Parse(conf.Template.Movie)
				tmplData := struct {
					Original media.Media
					Match    media.Movie
				}{m, matches[0]}

				var buf bytes.Buffer
				tmpl.Execute(&buf, tmplData)
				fmt.Printf("Rename `%s` to `%s`\n", m.Source(), buf.String())
				for i, candidate := range matches {
					buf.Reset()
					tmplData.Match = candidate
					tmpl.Execute(&buf, tmplData)
					fmt.Printf("  %d: %s\n", i, buf.String())
				}
			},
		},
	}
	app.Run(os.Args)
}

func loadConfig(ctx *cli.Context) config.Configuration {
	// TODO handle configuration load errors here
	conf, _ := config.Load(ctx.String("file"))

	return conf
}

func getMediaType(ctx *cli.Context, m media.Media) (media.MediaType, error) {
	switch ctx.String("type") {
	case "movie":
		return media.MOVIE, nil
	case "tvshow":
		return media.TVSHOW, nil
	case "auto":
		return media.GuessType(m), nil
	default:
		return media.UNKNOWN, fmt.Errorf(
			"`type` flag must be one of `movie`, `tvshow`, or `auto`; got %s\n",
			ctx.String("type"))
	}
}

func getQuery(ctx *cli.Context, m media.Media) string {
	if q := ctx.String("query"); len(q) > 0 {
		return q
	}
	return media.GuessName(m)
}
