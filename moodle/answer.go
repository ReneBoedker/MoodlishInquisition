package moodle

import (
	"fmt"
	"io"
)

// Answer describes a possible answer to a question.
// The object contains related information such as grading and feedback.
type Answer struct {
	text     string
	grade    float64
	feedback string
	options  map[string]string
}

// NewAnswer creates a new Answer object.
func NewAnswer(response string, grade float64) *Answer {
	return &Answer{
		text:     htmlEscape(response),
		grade:    grade,
		feedback: "",
		options:  make(map[string]string),
	}
}

// NewAnswerWithFeedback creates a new Answer object with specific feedback.
func NewAnswerWithFeedback(response string, grade float64, feedback string) *Answer {
	return &Answer{
		text:     htmlEscape(response),
		grade:    grade,
		feedback: htmlEscape(feedback),
		options:  make(map[string]string),
	}
}

// SetOption allows setting additional options for answers.
// For instance, one may use this to set 'tolerance' for numerical answers.
func (a *Answer) SetOption(option, value string) {
	a.options[option] = value
}

// GetOption retrieves the given option from a.
// If option has not been explicitly set, the return value isSet will be false.
func (a *Answer) GetOption(option string) (value string, isSet bool) {
	v, ok := a.options[option]
	return v, ok
}

// ToXml writes an Answer object to Moodle XML format.
// Note that this XML cannot be imported into Moodle on its own. It should be
// included in a QuestionBank to do so.
func (a *Answer) ToXml(w io.Writer) {
	fmt.Fprintf(w, `
	<answer fraction="%f">
		<text><![CDATA[%s]]></text>`, a.grade, a.text)
	defer fmt.Fprint(w, "\n\t</answer>")

	if a.feedback != "" {
		fmt.Fprintf(w, `
		<feedback format="html">
			<text><![CDATA[%s]]></text>
		</feedback>`, a.feedback)
	}

	for k, v := range a.options {
		fmt.Fprintf(w, `
		<%s>%s</%s>`, k, v, k)
	}
}
