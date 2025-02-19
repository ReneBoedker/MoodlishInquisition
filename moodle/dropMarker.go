package moodle

import (
	"fmt"
	"hash/fnv"
	"io"
	"strings"
)

var _ Question = (*DropMarker)(nil) // Ensure interface is satisfied

// Mark describes markers in the 'Drag and drop markers' question type.
type Mark struct {
	nDrags uint
	text   string
}

// NewMark creates a new marker to be dropped.
// If nDrags is zero, the number of copies will be set to infinite.
func NewMark(text string, nDrags uint) *Mark {
	return &Mark{
		text:   text,
		nDrags: nDrags,
	}
}

// Zone describes drop zones in the 'Drag and drop markers' question type.
type Zone struct {
	shape       string
	coords      string
	correctMark int
}

// NewZone defines a new drop zone with the given parameters.
// Supported shapes are 'circle' and 'rectangle'. The given coordinates describe
// the center of the zone. When defining a circle, its diameter will be
// max(width, height).
func NewZone(shape string, coords [2]float64, width, height float64, correctMark int) (*Zone, error) {
	var coordSpec string
	shape = strings.ToLower(shape)
	switch shape {
	case "circle":
		coordSpec = fmt.Sprintf("%.0f,%.0f;%.0f", coords[0], coords[1], max(width, height)/2)
	case "rectangle":
		coordSpec = fmt.Sprintf("%.0f,%.0f;%.0f,%.0f", coords[0]-width/2, coords[1]-height/2, width, height)
	default:
		return nil, fmt.Errorf("Unsupported shape %q", shape)
	}
	return &Zone{
		shape:       shape,
		coords:      coordSpec,
		correctMark: correctMark,
	}, nil
}

// DropMarker implements the 'Drag and drop marker' question type.
type DropMarker struct {
	name    string
	text    string
	file    string // base64-encoded
	points  uint
	shuffle bool
	markers []*Mark
	zones   []*Zone
}

// NewDropMarker creates a new 'Drag and drop marker' question.
// The 'file' argument must be a string containing the base64 encoded contents
// of the question image.
func NewDropMarker(description, file string, points uint, markers []*Mark, zones []*Zone) *DropMarker {
	hash := fnv.New32a()
	hash.Write([]byte(description))
	hash.Write([]byte(file))

	return &DropMarker{
		name:    fmt.Sprintf("%X", hash.Sum32()),
		text:    description,
		file:    file,
		points:  points,
		shuffle: true,
		markers: markers,
		zones:   zones,
	}
}

// MoodleName returns the question type as written in Moodle.
func (dm *DropMarker) MoodleName() string {
	return "Drag and drop markers"
}

// SetShuffleAnswers allows enabling or disabling shuffling of answers. The
// default is to shuffle.
func (dm *DropMarker) SetShuffleAnswers(b bool) {
	dm.shuffle = b
}

// ToXml writes a DropMarker object to Moodle XML format.
// Note that this XML cannot be imported into Moodle on its own. It should be
// included in a QuestionBank to do so.
func (dm *DropMarker) ToXml(w io.Writer) {
	// Write the question name and text
	fmt.Fprintf(w, `
<question type="ddmarker">
	<name>
		<text>%s</text>
	</name>
	<questiontext format="html">
		<text><![CDATA[`+"%s"+`]]></text>
	</questiontext>
	<defaultgrade>`+"%d"+`</defaultgrade>
	<showmisplaced/>
	<file name="figure.svg" encoding="base64">%s</file>`,
		dm.name, dm.text, dm.points, dm.file)
	defer fmt.Fprint(w, `
</question>`)

	if dm.shuffle {
		fmt.Fprintf(w, `
	<shuffleanswers>1</shuffleanswers>`)
	} else {
		fmt.Fprintf(w, `
	<shuffleanswers>0</shuffleanswers>`)
	}

	// Write the drag markers
	for i, v := range dm.markers {
		inf := ""
		if v.nDrags == 0 {
			// Add infinite specifier
			inf = `
		<infinite/>`
		}

		fmt.Fprintf(w, `
	<drag>
		<no>%d</no>
		<text>%s</text>%s
		<noofdrags>%d</noofdrags>
	</drag>`,
			i+1, v.text, inf, v.nDrags)
	}

	// Write the drop zones
	for i, v := range dm.zones {
		fmt.Fprintf(w, `
	<drop>
		<no>%d</no>
		<shape>%s</shape>
		<coords>%s</coords>
		<choice>%d</choice>
	</drop>`,
			i+1, v.shape, v.coords, v.correctMark+1)
	}
}
