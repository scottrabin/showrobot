package main

import (
	"bytes"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/scottrabin/showrobot/config"
	"github.com/scottrabin/showrobot/datasource"
	"github.com/scottrabin/showrobot/media"
	"github.com/scottrabin/showrobot/media/riff"
	"os"
	"path/filepath"
	"text/template"
)

var (
	fileFlag = cli.StringFlag{
		Name:  "file, f",
		Value: config.GetDefaultConfigurationPath(),
		Usage: "Use the specified configuration file",
	}
	typeFlag = cli.StringFlag{
		Name:  "type, t",
		Value: "auto",
		Usage: "Use the specified type for matching the media file",
	}
	queryFlag = cli.StringFlag{
		Name:  "query, q",
		Value: "",
		Usage: "Use the specified string as the search term instead of guessing from the file name",
	}
	interactiveFlag = cli.BoolFlag{
		Name:  "interactive, i",
		Usage: "Interactively list options to choose from",
	}
	noopFlag = cli.BoolFlag{
		Name:  "noop, n",
		Usage: "Report the intended action instead of performing it",
	}
)

type ResultFormat struct {
	Extension string
	Match     media.Movie
}

func main() {
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
			Name:        "rename",
			Usage:       "rename a file by finding information from the internet",
			Description: "Rename a file using information from an external service",
			Flags:       []cli.Flag{fileFlag, typeFlag, queryFlag, interactiveFlag, noopFlag},
			Action: func(c *cli.Context) {
				args := c.Args()
				conf := loadConfig(c)

				pwd, err := os.Getwd()
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				for _, p := range args {
					mf := media.NewFile(filepath.Join(pwd, p))

					mtype, err := getMediaType(c, mf)
					if err != nil {
						fmt.Println(err)
						os.Exit(1)
					}

					ds, err := datasource.New(conf, mtype)
					if err != nil {
						fmt.Println(err)
						os.Exit(1)
					}

					query := getQuery(c, mf)
					matches := ds.GetMovies(query)
					bestMatch := getBestMatch(c, mf, matches)

					// TODO handle errors here
					tmpl, _ := template.New("movie").Parse(conf.Template.Movie)

					var buf bytes.Buffer
					tmpl.Execute(&buf, ResultFormat{filepath.Ext(mf.Source), bestMatch})

					if c.Bool("noop") {
						// TODO Go doesn't have a way to escape shell arguments because the language
						// is actually well designed, but it would be nice to format the following anyway...
						fmt.Printf("mv \"%s\" \"%s\"\n", mf.Source, buf.String())
					} else {
						os.Rename(mf.Source, buf.String())
					}
				}
			},
		},
		{
			Name:        "riff",
			Usage:       "display RIFF data",
			Description: "Display RIFF data",
			Flags:       []cli.Flag{},
			Action: func(c *cli.Context) {
				args := c.Args()
				//conf := loadConfig(c)

				for _, p := range args {
					m := media.NewFile(p)

					riff.DoStuff(m.Source)
				}
			},
		},
	}
	app.Run(os.Args)
}

func loadConfig(ctx *cli.Context) config.Configuration {
	conf, err := config.Load(ctx.String("file"))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return conf
}

func getMediaType(ctx *cli.Context, mf *media.MediaFile) (media.MediaType, error) {
	switch ctx.String("type") {
	case "movie":
		return media.MOVIE, nil
	case "tvshow":
		return media.TVSHOW, nil
	case "auto":
		return media.GuessType(mf), nil
	default:
		return media.UNKNOWN, fmt.Errorf(
			"`type` flag must be one of `movie`, `tvshow`, or `auto`; got %s\n",
			ctx.String("type"))
	}
}

func getQuery(ctx *cli.Context, mf *media.MediaFile) string {
	if q := ctx.String("query"); len(q) > 0 {
		return q
	}
	return media.GuessName(mf)
}

func getBestMatch(ctx *cli.Context, mf *media.MediaFile, matches []media.Movie) media.Movie {
	if ctx.Bool("interactive") {
		fmt.Printf("Select match for [[ %s ]]\n", mf.Source)
		var optBuf bytes.Buffer
		optTmpl, _ := template.New("movieOption").Parse("{{.Name}} ({{.Year}})")
		for i, opt := range matches {
			optBuf.Reset()
			optTmpl.Execute(&optBuf, opt)
			fmt.Printf("  %d: %s\n", i+1, optBuf.String())
		}

		for {
			var choice int
			fmt.Printf("Enter choice: ")
			_, err := fmt.Scanf("%d", &choice)
			if err != nil {
				fmt.Println(err)
			}
			if 0 < choice && choice < len(matches)+1 {
				return matches[choice-1]
			}
		}
	}

	return matches[0]
}
