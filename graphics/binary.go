package graphics

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"
)

type binaryImage struct {
	content   []byte
	extension string
	alt       string
}

var _ Image = (*binaryImage)(nil) // Ensure that interface is satisfied

// ImageFromFile reads an image file into memory. An error is returned if
// filePath does not exist.
func ImageFromFile(path string) (*binaryImage, error) {
	content, err := fileAsBytes(path)

	ext := filepath.Ext(path)

	return &binaryImage{
		content:   content,
		extension: strings.TrimPrefix(ext, "."),
	}, err
}

// ImageFromBytes creates an image object directly from given byte slice.
// The slice should contain the same content as a file on disk of the specified
// file type. A simple sanity-check is performed, and an error is returned if
// the check fails. However, an image object with the specified content is
// produced regardless of the error value.
func ImageFromBytes(b []byte, filetype string) (*binaryImage, error) {
	var err error
	mime := http.DetectContentType(b)
	if !strings.HasPrefix(mime, "image/") {
		err = fmt.Errorf("Given bytes do not seem to be an image (detected MIME-type is %q)", mime)
	} else if ext := strings.TrimPrefix(mime, "image/"); ext != filetype {
		err = fmt.Errorf("Warning: File seems to be %s, not %s", ext, filetype)
	}

	return &binaryImage{
		content:   b,
		extension: strings.TrimPrefix(filetype, "."),
	}, err
}

// SetAltDescription allows given string to be used as the 'alt' attribute when
// converting to HTML.
func (img *binaryImage) SetAltDescription(s string) {
	img.alt = s
}

// Filetype returns the filetype of img.
func (img *binaryImage) Filetype() string {
	return img.extension
}

// ToHtml embeds img in Moodle-ready HTML and writes it to w.
// This should be used with care, especially with large image files, as they
// will be included directly in the HTML code.
func (img *binaryImage) ToHtml(w io.Writer) {
	fmt.Fprintf(w, `<img src="data:image/%s;base64,`, img.Filetype())

	img.ToBase64(w)

	if img.alt != "" {
		fmt.Fprintf(w, `" alt="%s" />`, img.alt)
	} else {
		fmt.Fprintf(w, `" />`)
	}
}

// ToBase64 encodes img to base64 format.
// This is for instance used to include graphics in the 'Drag and drop markers'
// question type.
func (img *binaryImage) ToBase64(w io.Writer) {
	fmt.Fprint(w, base64.StdEncoding.EncodeToString(img.content))
}
