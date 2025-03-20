package graphics

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
)

// cropSvg calls Inkscape to reduce the canvas size to its contents
func cropSvg(fileName string) error {
	version, err := inkscapeVersion()
	if err != nil {
		return err
	}

	var cmd *exec.Cmd
	// 'verb' command line arguments were removed in Inkscape 1.2
	if version[0] >= 1 && version[1] >= 2 {
		cmd = exec.Command("inkscape", `--actions="select-all;fit-canvas-to-selection;export-overwrite;export-do"`, fileName)
	} else {
		cmd = exec.Command("inkscape", "--verb=FitCanvasToDrawing", "--verb=FileSave", "--verb=FileQuit", fileName)
	}
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("inkscape: %v", err)
	}

	return nil
}

// inkscapeVersion extracts the version number of Inkscape (if installed)
func inkscapeVersion() ([]uint64, error) {
	cmd := exec.Command("inkscape", "--version")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	re := regexp.MustCompile(`^(?i:inkscape) ([.0-9]+)`)
	match := re.FindSubmatch(out)
	if len(match) == 0 {
		return nil, fmt.Errorf("Failed to extract version from %q", out)
	}

	versionBytes := bytes.Split(match[1], []byte("."))
	versionNums := make([]uint64, len(versionBytes))
	for i, v := range versionBytes {
		versionNums[i], err = strconv.ParseUint(string(v), 10, 0)
		if err != nil {
			return nil, err
		}
	}

	return versionNums, nil
}
