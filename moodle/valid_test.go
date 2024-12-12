package moodle

import (
	"encoding/xml"
	"io"
	"strings"
	"testing"
)

func exampleBank() *QuestionBank {
	q1 := NewMultiChoice(
		`$$f(t)=\int \cos(x)\, dx$$`,
		5,
		[]*Answer{
			NewAnswer("True", 100),
			NewAnswer("False", 0),
		},
	)
	q2 := NewMultiChoice(
		`$$x^2 + x = 0$$`,
		6,
		[]*Answer{
			NewAnswer("$$x=0$$", 100),
			NewAnswer("$$x=-1", 100),
			NewAnswer("$$x=4", 0),
		},
	)

	return &QuestionBank{
		name:      "Testing",
		questions: []Question{q1, q2},
	}
}

func TestOutputValid(t *testing.T) {
	qb := exampleBank()
	var b strings.Builder
	qb.ToXml(&b)

	d := xml.NewDecoder(strings.NewReader(b.String()))
	for {
		err := d.Decode(new(interface{}))
		if err != nil {
			if err == io.EOF {
				return
			}
			t.Fatalf("Decoding XML output produced error: %q", err)
		}
	}
}

func TestHtmlEscape(t *testing.T) {
	testCases := [...][2]string{
		{`$$(0,2\pi]$$`, `$$(0,2\pi&#93$$`},
		{`$$5<7$$`, `$$5\lt 7$$`},
		{`$$3>1$$`, `$$3\gt 1$$`},
	}
	for _, v := range testCases {
		if s := htmlEscape(v[0]); s != v[1] {
			t.Fatalf("Received %q, but expected %q", s, v[1])
		}
	}
}
