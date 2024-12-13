package moodle

import (
	"fmt"
	"hash/fnv"
	"io"
)

var _ Question = (*Numerical)(nil) // Ensure interface is satisfied

// MultiChoice implements the 'Multiple choice' question type.
type MultiChoice struct {
	name    string
	points  uint
	text    string
	answers []*Answer
}

// MoodleName returns the question type as written in Moodle.
func (q *MultiChoice) MoodleName() string {
	return `Multiple choice`
}

// NewMultiChoice creates a new 'Multiple choice' question.
func NewMultiChoice(description string, points uint, answers []*Answer) *MultiChoice {
	hash := fnv.New32a()
	hash.Write([]byte(description))
	for _, v := range answers {
		hash.Write([]byte(v.text))
	}

	return &MultiChoice{
		name:    fmt.Sprintf("%X", hash.Sum32()),
		points:  points,
		text:    description,
		answers: answers,
	}
}

// NCorrect counts the number of correct (incl. partially) answers in q.
func (q *MultiChoice) NCorrect() (n uint) {
	for _, a := range q.answers {
		if a.grade > 0 {
			n++
		}
	}
	return n
}

// GetDescription returns the description (i.e. the question text) of q.
func (q *MultiChoice) GetDescription() string {
	return q.text
}

// ToXml writes a MultiChoice object to Moodle XML format.
// Note that this XML cannot be imported into Moodle on its own. It should be
// included in a QuestionBank to do so.
func (q *MultiChoice) ToXml(w io.Writer) {
	// Write the question name and text
	fmt.Fprintf(w, `
<question type="multichoice">
	<name>
		<text>%s</text>
	</name>
	<questiontext format="html">
		<text><![CDATA[%s]]></text>
	</questiontext>
	<defaultgrade>%d</defaultgrade>`,
		q.name, q.text, q.points)
	defer fmt.Fprint(w, `
</question>`)

	// Several correct answers
	// One wrong answer will cancel out one correct answer
	for _, a := range q.answers {
		a.ToXml(w)
	}

	// Write remaining options
	fmt.Fprintf(w, "\n<shuffleanswers>1</shuffleanswers>")
	fmt.Fprintf(w, "\n<single>%t</single>", q.NCorrect() == 1)
	fmt.Fprintf(w, "\n<answernumbering>none</answernumbering>")
}
