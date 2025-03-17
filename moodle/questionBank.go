package moodle

import (
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
)

var reMathDelims = regexp.MustCompile(`(\$+)[^$]+(\$+)?`)

// QuestionBank is a collection of questions
type QuestionBank struct {
	name      string
	questions []Question
}

// NewQuestionBank creates a new question bank containing the given questions.
func NewQuestionBank(name string, questions []Question) *QuestionBank {
	return &QuestionBank{
		name:      name,
		questions: questions,
	}
}

// ToXml writes a QuestionBank object to Moodle XML format. The output is a
// complete file that can be imported in Moodle.
func (qb *QuestionBank) ToXml(w io.Writer) {
	fmt.Fprintf(w, `<?xml version="1.0" encoding="UTF-8"?>
<quiz>
<question type="category">
	<category>
		<text>$module$/%s</text>
	</category>
</question>`, qb.name)
	defer fmt.Fprint(w, `
</quiz>`)

	for _, q := range qb.questions {
		q.ToXml(w)
	}
}

// GenerateQuestionBank is the main function for generating random questions.
// It generates the specified number of questions and writes it to the given
// file (after creating it).
// An error is returned if the file already exists.
func GenerateQuestionBank(fName string, nQuestions int, gen func() Question) error {
	// Check that file does not exist
	if fileExists(fName) {
		return fmt.Errorf("File %q already exists", fName)
	}

	// Generate questions
	questions := make([]Question, nQuestions, nQuestions)
	for i := 0; i < nQuestions; i++ {
		questions[i] = gen()
		if err := validateSyntax(questions[i]); err != nil {
			return err
		}
	}
	qb := NewQuestionBank(fName, questions)

	f, err := os.Create(fName)
	if err != nil {
		return err
	}
	defer f.Close()
	qb.ToXml(f)
	return nil
}

func fileExists(fName string) bool {
	_, err := os.Stat(fName)
	return !errors.Is(err, os.ErrNotExist)
}

func validateSyntax(q Question) error {
	//return validateTeX(q.GetDescription())
	return nil
}

// validateTeX is currently not used
func validateTeX(s string) error {
	matches := reMathDelims.FindAllStringSubmatch(s, -1)
	for _, v := range matches {
		if len(v[1]) != 2 || len(v[2]) != 2 {
			return fmt.Errorf("Wrong or missing delimiters in string:\n%q", s)
		}
	}
	return nil
}
