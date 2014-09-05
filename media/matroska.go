package media

import (
	"fmt"
	"github.com/scottrabin/showrobot/media/format/ebml"
	"os"
	"strings"
	"syscall"
	"time"
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

	segment := m.findElement(els, ebml.ElementSegment)
	if segment == nil {
		return mi, fmt.Errorf("No segment element found")
	}
	segmentChildren, err := segment.Value()
	if err != nil {
		return mi, err
	}
	info := m.findElement(segmentChildren.([]ebml.Element), ebml.ElementInfo)
	if info == nil {
		return mi, fmt.Errorf("No info element found")
	}
	infoChildren, err := info.Value()
	if err != nil {
		return mi, err
	}
	timescale := ebml.Elements[ebml.ElementTimecodeScale].DefaultValue.(uint64)
	timescaleEl := m.findElement(infoChildren.([]ebml.Element), ebml.ElementTimecodeScale)
	if timescaleEl != nil {
		ts, _ := timescaleEl.Value()
		timescale = ts.(uint64)
	}
	durationEl := m.findElement(infoChildren.([]ebml.Element), ebml.ElementDuration)
	if durationEl == nil {
		return mi, fmt.Errorf("No duration element found")
	}
	duration, err := durationEl.Value()
	if err != nil {
		return mi, err
	}

	mi.Duration = time.Duration(uint64(duration.(float64)) * timescale)

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

func (m *matroskaCodec) findElement(els []ebml.Element, id ebml.EBMLID) *ebml.Element {
	for _, el := range els {
		if el.ID == id {
			return &el
		}
	}

	return nil
}

func init() {
	RegisterCodec(".mkv", &matroskaCodec{})
}
