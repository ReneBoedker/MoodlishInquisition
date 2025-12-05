package graphics

import (
	"bytes"
	_ "embed"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type SvgImage struct {
	content []byte
	dim     [2]float64
}

var _ Image = (*SvgImage)(nil) // Ensure that interface is satisfied

//go:embed preamble.tex
var preamble string

var svgDims = regexp.MustCompile(`((?:width|height)="([0-9.]+))(?:pt|px)?`)

// SvgFromFile reads an svg file into memory.
func SvgFromFile(path string) (*SvgImage, error) {
	// Read svg contents into memory
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Extract width and height
	w, h := 0.0, 0.0
	matches := svgDims.FindAllSubmatch(content, 2)
	for _, v := range matches {
		if bytes.HasPrefix(v[0], []byte("width")) {
			w, err = strconv.ParseFloat(string(v[2]), 0)
		} else {
			h, err = strconv.ParseFloat(string(v[2]), 0)
		}
	}
	if w == 0 || h == 0 || err != nil {
		return nil, fmt.Errorf("Failed to extract dimensions of svg.")
	}

	return &SvgImage{
		content: content,
		dim:     [2]float64{w, h},
	}, nil
}

// SvgFromTikz compiles a TikZ- or pfgplots-environment into an SvgImage.
// Intermediate results will be stored in tmpDir. If this argument is "", then
// a temporary folder will be created and deleted automatically. If tmpDir is
// specified, the called is responsible for deletion.
func SvgFromTikz(s string, tmpDir string) (*SvgImage, error) {
	svg, err := SvgFromMultipageTikz(s, tmpDir)
	if svg != nil {
		return svg[0], err
	}

	return nil, err
}

func SvgFromMultipageTikz(s string, tmpDir string) ([]*SvgImage, error) {
	if tmpDir == "" {
		var err error
		tmpDir, err = os.MkdirTemp("", "moodleTikz-*")
		if err != nil {
			return nil, err
		}
		defer os.RemoveAll(tmpDir)
	}

	svgPath, err := compileToSvg(s, tmpDir)
	if err != nil {
		return nil, err
	}

	svgs := make([]*SvgImage, len(svgPath), len(svgPath))
	for i, v := range svgPath {
		svgs[i], err = SvgFromFile(v)
		if err != nil {
			return nil, err
		}
	}

	return svgs, nil
}

// GetDimension returns the width and height of img as encoded in the svg file.
func (img *SvgImage) GetDimension() [2]float64 {
	return img.dim
}

// Scale changes the size of img by the given scaling factor. An error is
// returned if factor is not positive.
func (img *SvgImage) Scale(factor float64) error {
	if factor <= 0 {
		return fmt.Errorf("Scaling factor must be positive, but received %f", factor)
	}

	img.dim[0] *= factor
	img.dim[1] *= factor

	matches := svgDims.FindAllSubmatch(img.content, 2)
	for _, v := range matches {
		if bytes.HasPrefix(v[0], []byte("width")) {
			img.content = bytes.Replace(img.content, v[0], fmt.Appendf([]byte{}, `width="%.1fpx`, img.dim[0]), 1)
		} else {
			img.content = bytes.Replace(img.content, v[0], fmt.Appendf([]byte{}, `height="%.1fpx`, img.dim[1]), 1)
		}
	}

	return nil
}

// CropToContent will crop the image size to match the svg contents.
func (img *SvgImage) CropToContent() error {
	// Create temporary folder
	tmpDir, err := os.MkdirTemp("", "moodleTikz-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	// Write image to disk
	path := filepath.Join(tmpDir, "tmp.svg")
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	file.Write(img.content)
	file.Close()

	// Perform cropping
	err = cropSvg(path)
	if err != nil {
		return err
	}

	// Overwrite image contents with cropped image
	cropped, err := SvgFromFile(file.Name())
	if err != nil {
		return err
	}
	*img = *cropped

	return nil
}

// ToHtml embeds img in Moodle-ready HTML and writes it to w.
func (img *SvgImage) ToHtml(w io.Writer) {
	fmt.Fprintf(w, `<p>`)

	htmlContent := make([]byte, len(img.content))
	copy(htmlContent, img.content)

	// Remove the XML-tag
	htmlContent = regexp.MustCompile(`<\?xml.*?\?>\s*`).ReplaceAll(htmlContent, []byte(""))

	// Add Moodle's responsive image CSS-class
	htmlContent = regexp.MustCompile(`<svg`).ReplaceAll(htmlContent, []byte(`<svg class="img-responsive"`))

	fmt.Fprintf(w, "%s", htmlContent)
	fmt.Fprintf(w, "</p>\n")
}

// ToBase64 encodes img to base64 format.
// This is for instance used to include graphics in the 'Drag and drop markers'
// question type.
func (img *SvgImage) ToBase64(w io.Writer) {
	b64Content := make([]byte, len(img.content))
	copy(b64Content, img.content)

	// Change svg dimensions to px (to prevent bug in Moodle's implementation)
	b64Content = svgDims.ReplaceAll(b64Content, []byte(`${1}px`))

	fmt.Fprint(w, base64.StdEncoding.EncodeToString(b64Content))
}

// compileToPdf compiles a TikZ-picture into a PDF file.
// The output is the path of the resulting file.
func compileToPdf(s string, dir string) (string, error) {
	// Wrap tikzpicture in TeX-document
	var b strings.Builder
	fmt.Fprint(&b, preamble)
	fmt.Fprint(&b, s)
	fmt.Fprint(&b, "\n\\end{document}")

	// Compile file to pdf
	cmd := exec.Command("pdflatex", "--output-directory", dir, "--jobname", "tikz", "--")
	cmd.Stderr = os.Stderr
	cmd.Stdin = strings.NewReader(b.String())
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("pdflatex: %v", err)
	}

	return filepath.Join(dir, "tikz.pdf"), nil
}

// compileToSvg compiles a multipage TikZ-picture into individual SVG files.
// The output is a slice containing the path of each file.
func compileToSvg(s string, dir string) ([]string, error) {
	pdfPath, err := compileToPdf(s, dir)
	if err != nil {
		return nil, err
	}

	// Convert file to svg
	err = convertPdfToSvg(pdfPath, filepath.Join(dir, "tikz.svg"))
	if err != nil {
		return nil, err
	}

	return filepath.Glob(filepath.Join(dir, "tikz*.svg"))
}

// convertPdfToSvg will automatically call either pdftocairo or pdf2svg to convert given
// pdf file.
func convertPdfToSvg(pdfPath, destination string) error {
	err1 := pdftocairo(pdfPath, destination)
	if err1 == nil {
		// pdftocairo succeeded; no need to try pdf2svg
		return nil
	}

	err2 := pdf2svg(pdfPath, destination)
	if err2 != nil {
		return fmt.Errorf("%s\n%s", err1, err2)
	}

	return nil
}

// pdftocairo calls the external command of the same name to convert a PDF to SVG.
// It will convert every page into separate SVG files.
func pdftocairo(pdfPath, destination string) error {
	nPages, err := pdfPageCount(pdfPath)
	if err != nil {
		return fmt.Errorf(
			"pdfinfo failed, so pdftocairo cannot be used. Error message was: %s",
			err,
		)
	}

	for i := 1; i <= nPages; i++ {
		cmd := exec.Command(
			"pdftocairo",
			"-svg",
			"-f", fmt.Sprintf("%d", i),
			"-l", fmt.Sprintf("%d", i),
			pdfPath,
			strings.Replace(destination, ".svg", fmt.Sprintf("%02d.svg", i), 1))
		if err = cmd.Run(); err != nil {
			return fmt.Errorf("pdftocairo failed. Error message was: %s", err)
		}
	}

	return nil
}

// pdf2svg calls the external command of the same name to convert a PDF to SVG.
// It will convert every page into separate SVG files.
func pdf2svg(pdfPath, destination string) error {
	cmd := exec.Command(
		"pdf2svg",
		pdfPath,
		strings.Replace(destination, ".svg", "%02d.svg", 1),
		"all",
	)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("pdf2svg failed. Error message was: %s", err)
	}

	return nil
}

// pdfPageCount uses pdfinfo (part of poppler-utils) to extract page number.
func pdfPageCount(pdfPath string) (int, error) {
	var info strings.Builder
	cmd := exec.Command("pdfinfo", pdfPath)
	cmd.Stdout = &info
	if err := cmd.Run(); err != nil {
		return 0, err
	}

	pageStr := regexp.MustCompile(
		`Pages:[[:space:]]*([0-9]+)`,
	).FindStringSubmatch(info.String())

	if len(pageStr) == 0 {
		return 0, fmt.Errorf("Failed to extract number of PDF pages.")
	}
	i64, err := strconv.ParseInt(pageStr[1], 10, 0)
	if err != nil {
		return 0, err
	}
	return int(i64), nil
}

// Filetype returns the file type of img (i.e. "svg").
func (img *SvgImage) Filetype() string {
	return "svg"
}
