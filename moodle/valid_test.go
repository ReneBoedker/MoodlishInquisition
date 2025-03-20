package moodle

import (
	"encoding/xml"
	"io"
	"strings"
	"testing"
)

func exampleBank() *QuestionBank {
	questions := make([]Question, 0)

	questions = append(questions,
		NewMultiChoice(
			`$$f(t)=\int \cos(x)\, dx$$`,
			5,
			[]*Answer{
				NewAnswer("True", 100),
				NewAnswer("False", 0),
			},
		))

	questions = append(questions,
		NewDropText(
			`He is [[2]]`,
			1,
			[]*TextMark{
				NewTextMark("the Messiah", 0, false),
				NewTextMark("a very naughty boy", 0, true),
			},
		))

	questions = append(questions,
		NewShortText(
			`He is&hellip;`,
			1,
			[]*Answer{
				NewAnswer("the Messiah", 0),
				NewAnswer("a very naught boy", 100),
			},
		))

	piAns := NewAnswer("3.14", 100)
	piAns.SetOption("tolerance", "0")
	questions = append(questions,
		NewNumerical(
			`What is \(\pi\) rounded to two decimal places?`,
			1,
			[]*Answer{piAns},
		),
	)

	return &QuestionBank{
		name:      "Testing",
		questions: questions,
	}
}

func TestOutputValid(t *testing.T) {
	qb := exampleBank()
	var b strings.Builder
	qb.ToXml(&b)

	d := xml.NewDecoder(strings.NewReader(b.String()))
	for {
		err := d.Decode(new(any))
		if err != nil {
			if err == io.EOF {
				return
			}
			t.Fatalf("Decoding XML output produced error: %q", err)
		}
	}
}

func TestEscapeMath(t *testing.T) {
	testCases := [...][2]string{
		{`$$(0,2\pi]$$`, `$$(0,2\pi&#93$$`},
		{`$$5<7$$`, `$$5\lt 7$$`},
		{`$$3>1$$`, `$$3\gt 1$$`},
	}
	for _, v := range testCases {
		if s := EscapeMath(v[0]); s != v[1] {
			t.Fatalf("Received %q, but expected %q", s, v[1])
		}
	}
}
