package graphics

import (
	_ "embed"
	"encoding/xml"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

//go:embed example.tex
var example string

//go:embed exampleMulti.tex
var exampleMulti string

func TestSvgCompilation(t *testing.T) {
	img, err := SvgFromTikz(example, "")
	if err != nil {
		t.Fatalf("SvgFromTikz encountered error: %v", err)
	}

	err = img.Scale(2)
	if err != nil {
		t.Errorf("Scaling svg encountered error: %v", err)
	}

	err = img.CropToContent()
	if err != nil {
		t.Errorf("Cropping encountered error: %v", err)
	}
}

func TestMultipageSvg(t *testing.T) {
	imgs, err := SvgFromMultipageTikz(exampleMulti, "")
	if err != nil {
		t.Fatalf("SvgFromMultipageTikz encountered error: %v", err)
	}

	if len(imgs) != 2 {
		t.Errorf(
			"SvgFromMultipageTikz produced %d images, but 2 were expected",
			len(imgs),
		)
	}
}

// Test pdf2svg manually, as pdftocairo takes priority
func TestPdf2svg(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "moodleTikz-*")
	if err != nil {
		t.Fatalf("Failed to create temporary folder: %s", err)
	}

	path, err := compileToPdf(exampleMulti, tmpDir)
	err = pdf2svg(path, filepath.Join(tmpDir, "tikz.svg"))

	generatedFiles, _ := filepath.Glob(filepath.Join(tmpDir, "tikz*.svg"))
	if len(generatedFiles) != 2 {
		t.Errorf(
			"pdf2svg produced %d images, but 2 were expected",
			len(generatedFiles),
		)
	}
}

// Test that HTML output is valid HTML
func TestHtmlValid(t *testing.T) {
	img, err := SvgFromTikz(example, "")
	if err != nil {
		t.Fatalf("SvgFromTikz encountered error: %v", err)
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
