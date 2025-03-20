package graphics

import (
	"bytes"
	_ "embed"
	"encoding/xml"
	"io"
	"testing"
)

//go:embed everest.png
var pngExample []byte

//go:embed everest.jpeg
var jpegExample []byte

func TestBinaryInputs(t *testing.T) {
	tests := []*binaryImage{
		{content: pngExample, extension: "png"},
		{content: jpegExample, extension: "jpeg"},
	}

	for i, v := range tests {
		img, err := ImageFromBytes(v.content, v.extension)
		if err != nil {
			t.Fatalf("Creating image produced error: %s", err)
		}

		if i == 0 {
			img.SetAltDescription("This is a test")
		}

		var b bytes.Buffer
		img.ToHtml(&b)

		d := xml.NewDecoder(&b)
		for {
			err := d.Decode(new(any))
			if err != nil {
				if err == io.EOF {
					break
				}
				t.Fatalf("Decoding XML output produced error: %s", err)
			}
		}
	}
}
