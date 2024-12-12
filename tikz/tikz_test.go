package tikz

import (
	_ "embed"
	"testing"
)

//go:embed example.tex
var example string

func TestCompile(t *testing.T) {
	_, err := CompileToHtml(example, false, "")
	if err != nil {
		t.Fatalf("CompileToHtml encountered error: %v", err)
	}

	_, err = CompileToHtml(example, true, "")
	if err != nil {
		t.Logf("inkscape encountered error: %v", err)
	}

	_, _, err = CompileToBase64(example, 1, false, "")
	if err != nil {
		t.Fatalf("CompileToBase64 encountered error: %v", err)
	}
}
