package config

import "encoding/json"
import "fmt"
import "io/ioutil"
import "os"
import "os/user"
import "path/filepath"
import "strings"

type Configuration struct {
	Format   string             `json:"format"`
	Language string             `json:"language"`
	ApiKey   map[string]*string `json:"apikey"`
	Template struct {
		Movie  string
		TVShow string
	}
}

type InvalidPropertyAccess struct {
	Property string
}

func (err *InvalidPropertyAccess) Error() string {
	return fmt.Sprintf("Invalid configuration property: `%s`", err.Property)
}

func GetDefaultConfigurationPath() string {
	usr, err := user.Current()
	if err != nil {
		return ""
	}

	return filepath.Join(usr.HomeDir, ".showrobot/config.json")
}

func Load(file string) (conf Configuration, err error) {
	conf.Format = ""
	conf.Language = "en"
	conf.ApiKey = make(map[string]*string)
	conf.Template.Movie = "{{.Match.Name}} ({{.Match.Year}}){{.Original.GetExtension}}"
	conf.Template.TVShow = "TODO"

	contents, err := ioutil.ReadFile(file)
	if err == nil {
		json.Unmarshal(contents, &conf)
	}

	return
}

func (conf *Configuration) Save(file string) error {
	bytes, err := json.Marshal(conf)
	if err != nil {
		return err
	}
	err = os.MkdirAll(filepath.Dir(file), os.ModePerm)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(file, bytes, os.ModePerm)

	return err
}

func (conf *Configuration) getConfigProperty(property string) (*string, error) {
	var err error

	fields := strings.Split(property, ".")

	switch fields[0] {
	case "format":
		return &conf.Format, nil
	case "language":
		return &conf.Language, nil
	case "apikey":
		if len(fields) == 2 {
			valueAddr, has := conf.ApiKey[fields[1]]
			if !has {
				valueAddr = new(string)
				conf.ApiKey[fields[1]] = valueAddr
			}
			return valueAddr, nil
		} else {
			err = fmt.Errorf("Setting an API key requires an API subkey")
		}
	case "template":
		if len(fields) == 2 {
			switch fields[1] {
			case "movie":
				return &conf.Template.Movie, nil
			case "tvshow":
				return &conf.Template.TVShow, nil
			}
		}
		err = fmt.Errorf("Setting a template requires a subkey (`movie` or `tvshow`)")
	default:
		err = &InvalidPropertyAccess{fields[0]}
	}

	return nil, err
}

func (conf *Configuration) Get(property string) (value string, err error) {
	valueAddr, err := conf.getConfigProperty(property)

	if err == nil {
		value = *valueAddr
	}

	return
}

func (conf *Configuration) Set(property string, value string) error {
	valueAddr, err := conf.getConfigProperty(property)

	if err == nil {
		*valueAddr = value
	}

	return err
}
