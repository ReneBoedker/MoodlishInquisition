package moodle

import "strings"

// EscapeMath ensures that LaTeX math is not interpreted as HTML code when
// imported into Moodle.
func EscapeMath(s string) string {
	replacements := [...][2]string{
		{`[`, `&#91`},
		{`]`, `&#93`},
		{`<`, `\lt `},
		{`>`, `\gt `},
	}

	for _, v := range replacements {
		s = strings.ReplaceAll(s, v[0], v[1])
	}
	return s
}
