package moodle

import (
	"fmt"
	"hash/fnv"
	"io"
)

var _ Question = (*DropText)(nil) // Ensure interface is satisfied

// TextMark describes markers in the 'Drag and drop into text' question type.
type TextMark struct {
	text      string
	dropGroup uint
	unlimited bool
}

// NewTextMark creates a new text marker to be dropped.
func NewTextMark(text string, group uint, unlimited bool) *TextMark {
	return &TextMark{
		text:      text,
		dropGroup: group,
		unlimited: unlimited,
	}
}

// DropText implements the 'Drag and drop into text' question type.
type DropText struct {
	name    string
	text    string
	points  uint
	shuffle bool
	markers []*TextMark
}

// NewDropText creates a new 'Drag and drop text' question.
// The description should contain substrings of the form [[n]], where n
// corresponds to the index of one of the markers. Note that indexing in the
// description must match Moodle's internal indexing, i.e. starting from 1.
// Thus, mark [[n]] in the description matches marker n-1 in the slice of
// markers.
func NewDropText(description string, points uint, markers []*TextMark) *DropText {
	hash := fnv.New32a()
	hash.Write([]byte(description))

	return &DropText{
		name:    fmt.Sprintf("%X", hash.Sum32()),
		text:    description,
		points:  points,
		shuffle: true,
		markers: markers,
	}
}

// MoodleName returns the question type as written in Moodle.
func (dt *DropText) MoodleName() string {
	return "Drag and drop into text"
}

// SetShuffleAnswers allows enabling or disabling shuffling of answers. The
// default is to shuffle.
func (dt *DropText) SetShuffleAnswers(b bool) {
	dt.shuffle = b
}

// ToXml writes a DropText object to Moodle XML format.
// Note that this XML cannot be imported into Moodle on its own. It should be
// included in a QuestionBank to do so.
func (dt *DropText) ToXml(w io.Writer) {
	// Write the question name and text
	fmt.Fprintf(w, `
<question type="ddwtos">
	<name>
		<text>%s</text>
	</name>
	<questiontext format="html">
		<text><![CDATA[`+"%s"+`]]></text>
	</questiontext>
	<defaultgrade>`+"%d"+`</defaultgrade>`,
		dt.name, dt.text, dt.points)
	defer fmt.Fprint(w, `
</question>`)

	if dt.shuffle {
		fmt.Fprintf(w, `
	<shuffleanswers/>`)
	}

	// Write markers
	for _, v := range dt.markers {
		inf := ""
		if v.unlimited {
			// Add infinite specifier
			inf = `
		<infinite/>`
		}

		fmt.Fprintf(w, `
	<dragbox>
		<text>%s</text>
		<group>%d</group>%s
	</dragbox>`,
			v.text, v.dropGroup+1, inf)
	}
}
