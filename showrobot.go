package main

import "fmt"
import "github.com/codegangsta/cli"
import "github.com/scottrabin/showrobot/config"
import "os"

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
			Flags: []cli.Flag{
				cli.StringFlag{"file", config.GetDefaultConfigurationPath(), "Use the specified configuration file"},
			},
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
	}
	app.Run(os.Args)
}
