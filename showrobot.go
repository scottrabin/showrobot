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
				conf, _ := config.Load(c.String("file"))

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
			Flags:       []cli.Flag{fileFlag},
			Action: func(c *cli.Context) {
				args := c.Args()
				conf, _ := config.Load(c.String("file"))

				mf, err := media.NewMedia(args[0])
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				ds := datasource.NewMovieSource(conf, "themoviedb")
				movies := ds.GetMovies(mf)

				// TODO sort movies by match likelihood

				// TODO handle errors here
				tmpl, _ := template.New("movie").Parse(conf.Template.Movie)
				tmplData := struct {
					Original media.Media
					Match    media.Movie
				}{mf, movies[0]}

				var buf bytes.Buffer
				tmpl.Execute(&buf, tmplData)
				fmt.Printf("Rename `%s` to `%s`\n", mf.GetFileName(), buf.String())
				for i, candidate := range movies {
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
