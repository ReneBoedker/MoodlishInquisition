package moodle

import "strings"

func htmlEscape(s string) string {
	if strings.Contains(s, "svg") || strings.Contains(s, "<p>") {
		return s
	}

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
