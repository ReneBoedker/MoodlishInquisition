package graphics

import (
	_ "embed"
	"encoding/xml"
	"io"
	"strings"
	"testing"
)

//go:embed example.tex
var example string

func TestSvgCompilation(t *testing.T) {
	img, err := SvgFromTikz(example, false, "")
	if err != nil {
		t.Fatalf("SvgFromTikz (without cropping) encountered error: %v", err)
	}

	err = img.Scale(2)
	if err != nil {
		t.Fatalf("Scaling svg encountered error: %v", err)
	}

	_, err = SvgFromTikz(example, true, "")
	if err != nil {
		t.Logf("SvgFromTikz (with cropping) encountered error: %v", err)
	}
}

// Test that HTML output is valid HTML
func TestHtmlValid(t *testing.T) {
	img, err := SvgFromTikz(example, false, "")
	if err != nil {
		t.Fatalf("SvgFromTikz (without cropping) encountered error: %v", err)
	}

	var b strings.Builder
	img.ToHtml(&b)

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
