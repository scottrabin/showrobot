package media

import "time"

type aviMediaFormat struct{}

func (media *aviMediaFormat) Duration(m Media) time.Duration {
	// TODO find a way to determine the actual runtime of the file
	return time.Duration(90) * time.Minute
}

func init() {
	Register(".avi", &aviMediaFormat{})
}
