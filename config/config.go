package config

import "encoding/json"
import "fmt"
import "io/ioutil"
import "os"
import "os/user"
import "path/filepath"

type Configuration struct {
	Format   string `json:"format"`
	Language string `json:"language"`
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

func Load(file string) (Configuration, error) {
	conf := Configuration{
		Format:   "",
		Language: "en",
	}

	contents, err := ioutil.ReadFile(file)
	if err == nil {
		json.Unmarshal(contents, &conf)
	}

	return conf, err
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

func (conf *Configuration) Get(property string) (value string, err error) {
	switch property {
	case "format":
		value = conf.Format
	case "language":
		value = conf.Language
	default:
		err = &InvalidPropertyAccess{property}
	}

	return
}

func (conf *Configuration) Set(property string, value string) (err error) {
	switch property {
	case "format":
		conf.Format = value
	case "language":
		conf.Language = value
	default:
		err = &InvalidPropertyAccess{property}
	}

	return
}
