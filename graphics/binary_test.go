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
			t.Errorf("Creating image produced error: %s", err)
			continue
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
				t.Errorf("Decoding XML output produced error: %s", err)
			}
		}
	}
}

func TestMismatchedInputs(t *testing.T) {
	_, err := ImageFromBytes([]byte("test"), "gif")
	if err == nil {
		t.Errorf("Non-image file failed to return an error")
	}

	_, err = ImageFromBytes(pngExample, "jpeg")
	if err == nil {
		t.Errorf("Mismatched file type failed to return an error")
	}
}
