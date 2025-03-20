package graphics

import (
	"fmt"
	"os"
)

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
