package graphics

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFileHandling(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "moodlishInquisitionTest-*")
	if err != nil {
		t.Fatalf("Creating temporary folder caused error: %s", err)
	}
	defer os.RemoveAll(tmpDir)

	path := filepath.Join(tmpDir, "test.txt")
	if pathExists(path) {
		t.Fatalf("pathExists returned true for non-existent file")
	}

	file, err := os.Create(path)
	if err != nil {
		t.Fatalf("Creating test file caused error: %s", err)
	}
	testString := "This is a test"
	_, err = file.WriteString(testString)
	if err != nil {
		t.Fatalf("Writing to test file caused error: %s", err)
	}
	file.Close()

	if !pathExists(path) {
		t.Fatalf("pathExists returned false for existent file")
	}

	content, err := fileAsBytes(path)
	if err != nil {
		t.Fatalf("Reading from file caused error: %s", err)
	}
	if string(content) != testString {
		t.Fatalf("Expected file contents %q, but got %q", testString, content)
	}
}
