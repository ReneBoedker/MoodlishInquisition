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
		text:     response,
		grade:    grade,
		feedback: "",
		options:  make(map[string]string),
	}
}

// NewAnswerWithFeedback creates a new Answer object with specific feedback.
func NewAnswerWithFeedback(response string, grade float64, feedback string) *Answer {
	return &Answer{
		text:     response,
		grade:    grade,
		feedback: feedback,
		options:  make(map[string]string),
	}
}

// SetGrade sets the answer grade to the given value. It returns an error if
// grade is not in the interval [-100, 100].
func (a *Answer) SetGrade(grade float64) error {
	if grade < -100 || grade > 100 {
		return fmt.Errorf("Grade must be between")
	}
	a.grade = grade
	return nil
}

// GetGrade returns the current grade of the answer.
func (a *Answer) GetGrade() float64 {
	return a.grade
}

// SetOption allows setting additional options for answers.
// For instance, one may use this to set 'tolerance' for numerical answers.
//
// The function will not check if the specified option and its value are valid.
func (a *Answer) SetOption(option, value string) {
	a.options[option] = value
}

// GetOption retrieves the given option from a.
// If option has not been explicitly set, the return value isSet will be false.
func (a *Answer) GetOption(option string) (value string, isSet bool) {
	v, ok := a.options[option]
	return v, ok
}

// SetFeedback overwrites existing feedback with given string.
func (a *Answer) SetFeedback(s string) {
	a.feedback = s
}

// GetFeedback returns the current feedback of a.
func (a *Answer) GetFeedback() string {
	return a.feedback
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
