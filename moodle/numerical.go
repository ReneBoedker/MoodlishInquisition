package moodle

import (
	"fmt"
	"hash/fnv"
	"io"
)

var _ Question = (*Numerical)(nil) // Ensure interface is satisfied

// Numerical implements the 'Numerical' question type in Moodle
type Numerical struct {
	name    string
	points  uint
	text    string
	answers []*Answer
}

// MoodleName returns the question type as written in Moodle.
func (q *Numerical) MoodleName() string {
	return `Numerical`
}

// NewNumerical creates a new 'Numerical' question.
func NewNumerical(description string, points uint, answers []*Answer) *Numerical {
	hash := fnv.New32a()
	hash.Write([]byte(description))
	for _, v := range answers {
		hash.Write([]byte(v.text))
	}

	return &Numerical{
		name:    fmt.Sprintf("%X", hash.Sum32()),
		points:  points,
		text:    description,
		answers: answers,
	}
}

// SetShuffleAnswers allows enabling or disabling shuffling of answers. This has
// no effect for Numerical question types.
func (q *Numerical) SetShuffleAnswers(b bool) {
}

// ToXml writes a Numerical object to Moodle XML format.
// Note that this XML cannot be imported into Moodle on its own. It should be
// included in a QuestionBank to do so.
func (q *Numerical) ToXml(w io.Writer) {
	// Write the question name and text
	fmt.Fprintf(w, `
<question type="numerical">
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
	fmt.Fprint(w, "\n<unitgradingtype>0</unitgradingtype>")
}
