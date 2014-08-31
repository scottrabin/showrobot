package media

import "time"

var AVI = Codec{
	Decode: func(mf *MediaFile) (MediaInfo, error) {
		mi := MediaInfo{
			// TODO find a way to determine the actual runtime of the file
			Duration: 90 * time.Minute,
		}

		return mi, nil
	},
}

func init() {
	RegisterCodec(".avi", &AVI)
}
