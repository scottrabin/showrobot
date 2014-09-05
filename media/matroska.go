package media

import (
	"fmt"
	"github.com/scottrabin/showrobot/media/format/ebml"
	"os"
	"strings"
	"syscall"
)

type matroskaCodec struct{}

func (m *matroskaCodec) Decode(mf *MediaFile) (MediaInfo, error) {
	mi := MediaInfo{}

	file, err := os.Open(mf.Source)
	if err != nil {
		return mi, err
	}
	fileinfo, err := file.Stat()
	if err != nil {
		return mi, err
	}

	mmap, err := syscall.Mmap(int(file.Fd()), 0, int(fileinfo.Size()),
		syscall.PROT_READ, syscall.MAP_SHARED)
	if err != nil {
		return mi, err
	}

	els, _ := ebml.ReadElements(mmap)
	m.PrintElements(els, 0)

	return mi, nil
}

func (m *matroskaCodec) PrintElements(els []ebml.Element, depth int) {
	for _, el := range els {
		v, err := el.Value()
		if err != nil {
			fmt.Printf("%s%x (%s) => Error: %v\n",
				strings.Repeat(" ", depth*4), el.ID, el.Meta().Name, err)
			continue
		}
		switch v := v.(type) {
		case []ebml.Element:
			fmt.Printf("%s%x (%s)\n", strings.Repeat(" ", depth*4), el.ID, el.Meta().Name)
			if el.ID != ebml.ElementCluster {
				// skip clusters; for everything else, show children
				m.PrintElements(v, depth+1)
			}
		case []byte:
			if len(v) > 24 {
				v = v[:24]
			}
			fmt.Printf("%s%x (%s)\t%v\n", strings.Repeat(" ", depth*4), el.ID, el.Meta().Name, v)
		default:
			fmt.Printf("%s%x (%s)\t%v\n",
				strings.Repeat(" ", depth*4), el.ID, el.Meta().Name, v)
		}
	}
}

func init() {
	RegisterCodec(".mkv", &matroskaCodec{})
}
