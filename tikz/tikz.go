package tikz

import (
	"bytes"
	_ "embed"
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

//go:embed preamble.tex
var preamble string

var svgDims = regexp.MustCompile(`((?:width|height)="([0-9.]+))(?:pt|px)`)

// CompileToHtml produces an HTML-version of a tikzpicture-environment.
// If argument crop is true, Inkscape will be used to crop the figure to its
// contents. Intermediate results will be stored in tmpDir. If this argument is
// "", then a temporary folder will be created and deleted automatically. If
// tmpDir is specified, the caller is responsible for deletion.
func CompileToHtml(s string, crop bool, tmpDir string) (string, error) {
	if tmpDir == "" {
		var err error
		tmpDir, err = os.MkdirTemp("", "moodleTikz-*")
		if err != nil {
			return "", err
		}
		defer os.RemoveAll(tmpDir)
	}

	svgPath, err := compileToSvg(s, crop, tmpDir)
	if err != nil {
		return "", err
	}

	out, err := svgToHtml(svgPath)
	if err != nil {
		return "", fmt.Errorf("svgToHtml: %v", err)
	}

	return out, nil
}

// CompileToBase64 compiles a LaTeX-string, and produces a base64 encoding of
// the resulting svg.
// If argument crop is true, Inkscape will be used to crop the figure to its
// contents. Intermediate results will be stored in tmpDir. If this argument is
// "", then a temporary folder will be created and deleted automatically. If
// tmpDir is specified, the caller is responsible for deletion.
func CompileToBase64(s string, scale float64, crop bool, tmpDir string) (encoded string, dim [2]float64, err error) {
	if tmpDir == "" {
		var err error
		tmpDir, err = os.MkdirTemp("", "moodleTikz-*")
		if err != nil {
			return "", [2]float64{0, 0}, err
		}
		defer os.RemoveAll(tmpDir)
	}

	svgPath, err := compileToSvg(s, crop, tmpDir)
	if err != nil {
		return "", [2]float64{0, 0}, err
	}

	// Read svg contents into memory
	file, err := os.ReadFile(svgPath)
	if err != nil {
		return "", [2]float64{0, 0}, err
	}

	// Change svg dimensions to px (to prevent bug in Moodle's implementation)
	file = svgDims.ReplaceAll(file, []byte(`${1}px`))

	// Extract width and height
	w, h := 0.0, 0.0
	matches := svgDims.FindAllSubmatch(file, 2)
	for _, v := range matches {
		if bytes.HasPrefix(v[0], []byte("width")) {
			w, err = strconv.ParseFloat(string(v[2]), 0)
			file = bytes.Replace(file, v[0], []byte(fmt.Sprintf(`width="%.1fpx`, scale*w)), 1)
		} else {
			h, err = strconv.ParseFloat(string(v[2]), 0)
			file = bytes.Replace(file, v[0], []byte(fmt.Sprintf(`height="%.1fpx`, scale*h)), 1)
		}
	}
	if w == 0 || h == 0 || err != nil {
		return "", [2]float64{0, 0}, fmt.Errorf("Failed to extract dimensions of svg.")
	}

	// Encode as base64
	return base64.StdEncoding.EncodeToString(file), [2]float64{w * scale, h * scale}, nil
}

func compileToSvg(s string, crop bool, dir string) (string, error) {
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

	// Convert file to svg
	cmd = exec.Command("pdf2svg", filepath.Join(dir, "tikz.pdf"), filepath.Join(dir, "tikz.svg"))
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("pdf2svg: %v", err)
	}

	// Crop unnecessary space around figure if requested
	if crop {
		cmd = exec.Command("inkscape", "--verb=FitCanvasToDrawing", "--verb=FileSave", "--verb=FileQuit", filepath.Join(dir, "tikz.svg"))
		if err := cmd.Run(); err != nil {
			return "", fmt.Errorf("inkscape: %v", err)
		}
	}

	return filepath.Join(dir, "tikz.svg"), nil
}

// svgToHtml reads an SVG file from disk, and returns Moodle-ready content.
func svgToHtml(path string) (string, error) {
	var b strings.Builder
	fmt.Fprint(&b, "<p>")

	// Read SVG-file and remove the XML-tag
	imgBytes, err := fileAsBytes(path)
	if err != nil {
		return "", err
	}
	imgBytes = regexp.MustCompile(`<\?xml.*?\?>\s*`).ReplaceAll(imgBytes, []byte(""))

	// Add Moodle's responsive image CSS-class
	imgBytes = regexp.MustCompile(`<svg`).ReplaceAll(imgBytes, []byte(`<svg class="img-responsive"`))

	fmt.Fprintf(&b, "%s", imgBytes)
	fmt.Fprint(&b, "</p>\n<p>")
	return b.String(), nil
}

// fileAsBytes reads an entire file into memory, and returns it as a byte slice.
func fileAsBytes(path string) ([]byte, error) {
	if !pathExists(path) {
		return nil, fmt.Errorf("Requested file %q does not exist.", path)
	}
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return file, nil
}

// pathExists checks if a given file already exists.
func pathExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}
