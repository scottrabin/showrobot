package media

import "time"

type AviMediaFile struct {
	filename string
	runtime  time.Duration
}

func (media *AviMediaFile) GetRuntime() time.Duration {
	// TODO find a way to determine the actual runtime of the file
	return time.Duration(90) * time.Minute
}

func (media *AviMediaFile) GetName() string {
	return media.filename
}
