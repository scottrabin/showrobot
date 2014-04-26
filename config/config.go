package config

import "fmt"

type Configuration struct {
  Format string
  Language string
}

type InvalidPropertyAccess struct {
  Property string
}

func (err *InvalidPropertyAccess) Error() string {
  return fmt.Sprintf("Invalid configuration property: `%s`", err.Property)
}

func Load() (conf Configuration) {
  // TODO load from base file and from user settings
  conf.Format = "{n}"
  conf.Language = "en"

  return
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

func (conf *Configuration) Save() (err error) {
  // TODO
  return
}
