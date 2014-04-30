package media

import "time"

type AviMediaFile struct {
	MediaFile
}

func (media *AviMediaFile) GetRuntime() time.Duration {
	// TODO find a way to determine the actual runtime of the file
	return time.Duration(90) * time.Minute
}
