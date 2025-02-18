package moodle

import (
	"fmt"
	"hash/fnv"
	"io"
)

var _ Question = (*ShortText)(nil) // Ensure interface is satisfied

// ShortText implements the 'Short-Answer' question type.
type ShortText struct {
	name          string
	points        uint
	text          string
	answers       []*Answer
	caseSensitive bool
}

// MoodleName returns the question type as written in Moodle.
func (q *ShortText) MoodleName() string {
	return `Short Answer`
}

// NewShortText creates a new 'Short-Answer' question.
func NewShortText(description string, points uint, answers []*Answer) *ShortText {
	hash := fnv.New32a()
	hash.Write([]byte(description))
	for _, v := range answers {
		hash.Write([]byte(v.text))
	}

	return &ShortText{
		name:          fmt.Sprintf("%X", hash.Sum32()),
		points:        points,
		text:          description,
		answers:       answers,
		caseSensitive: false,
	}
}

// SetCaseSensitivity sets the case sensitivity of q.
// The default for a new question is false.
func (q *ShortText) SetCaseSensitivity(b bool) {
	q.caseSensitive = b
}

// SetShuffleAnswers allows enabling or disabling shuffling of answers. This has
// no effect for 'Short Answer' question types.
func (q *ShortText) SetShuffleAnswers(b bool) {
}

// ToXml writes a ShortText object to Moodle XML format.
// Note that this XML cannot be imported into Moodle on its own. It should be
// included in a QuestionBank to do so.
func (q *ShortText) ToXml(w io.Writer) {
	// Write the question name and text
	fmt.Fprintf(w, `
<question type="shortanswer">
	<name>
		<text>%s</text>
	</name>
	<questiontext format="html">
		<text><![CDATA[%s]]></text>
	</questiontext>
	<defaultgrade>%d</defaultgrade>`,
		q.name, q.text, q.points)
	defer fmt.Fprintf(w, `
</question>`)

	// Write answers
	for _, a := range q.answers {
		a.ToXml(w)
	}

	// Write remaining options
	if q.caseSensitive {
		fmt.Fprint(w, "\n<usecase>1</usecase>")
	} else {
		fmt.Fprint(w, "\n<usecase>0</usecase>")
	}
}
