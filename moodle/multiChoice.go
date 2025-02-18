package moodle

import (
	"fmt"
	"hash/fnv"
	"io"
)

var _ Question = (*MultiChoice)(nil) // Ensure interface is satisfied

// MultiChoice implements the 'Multiple choice' question type.
type MultiChoice struct {
	name    string
	points  uint
	shuffle bool
	text    string
	answers []*Answer
}

// MoodleName returns the question type as written in Moodle.
func (mc *MultiChoice) MoodleName() string {
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
		shuffle: true,
		text:    description,
		answers: answers,
	}
}

// NCorrect counts the number of correct (incl. partially) answers in q.
func (mc *MultiChoice) NCorrect() (n uint) {
	for _, a := range mc.answers {
		if a.grade > 0 {
			n++
		}
	}
	return n
}

// GetDescription returns the description (i.e. the question text) of q.
func (mc *MultiChoice) GetDescription() string {
	return mc.text
}

// SetShuffleAnswers allows enabling or disabling shuffling of answers. The
// default is to shuffle.
func (mc *MultiChoice) SetShuffleAnswers(b bool) {
	mc.shuffle = b
}

// ToXml writes a MultiChoice object to Moodle XML format.
// Note that this XML cannot be imported into Moodle on its own. It should be
// included in a QuestionBank to do so.
func (mc *MultiChoice) ToXml(w io.Writer) {
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
		mc.name, mc.text, mc.points)
	defer fmt.Fprint(w, `
</question>`)

	if mc.shuffle {
		fmt.Fprintf(w, `
	<shuffleanswers/>`)
	}

	// Several correct answers
	// One wrong answer will cancel out one correct answer
	for _, a := range mc.answers {
		a.ToXml(w)
	}

	// Write remaining options
	fmt.Fprintf(w, "\n<shuffleanswers>1</shuffleanswers>")
	fmt.Fprintf(w, "\n<single>%t</single>", mc.NCorrect() == 1)
	fmt.Fprintf(w, "\n<answernumbering>none</answernumbering>")
}
