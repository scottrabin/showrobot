package media

import "time"

type unknownMediaFormat struct{}

func (media *unknownMediaFormat) Duration(m Media) time.Duration {
	return time.Duration(0)
}

func init() {
	Register("", &unknownMediaFormat{})
}
