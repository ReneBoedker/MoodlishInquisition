package graphics

import (
	"encoding/base64"
	"fmt"
	"io"
	"path/filepath"
	"strings"
)

type binaryImage struct {
	content   []byte
	extension string
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
// file type.
// Note that the function will not perform any checks to ensure that the
// encoding is valid.
func ImageFromBytes(b []byte, filetype string) *binaryImage {
	return &binaryImage{
		content:   b,
		extension: strings.TrimPrefix(filetype, "."),
	}
}

// Filetype returns the filetype of img.
func (img *binaryImage) Filetype() string {
	return img.extension
}

// ToHtml embeds img in Moodle-ready HTML and writes it to w.
// This should be used with care, especially with large image files, as they
// will be included directly in the HTML code.
func (img *binaryImage) ToHtml(w io.Writer) {
	fmt.Fprintf(w, `<img src="data:image/%s;base64,"`, img.Filetype())
	img.ToBase64(w)
	fmt.Fprintf(w, `">"`)
}

// ToBase64 encodes img to base64 format.
// This is for instance used to include graphics in the 'Drag and drop markers'
// question type.
func (img *binaryImage) ToBase64(w io.Writer) {
	fmt.Fprint(w, base64.StdEncoding.EncodeToString(img.content))
}
