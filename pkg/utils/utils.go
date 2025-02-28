package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// OpenFileForReading opens a file for reading
func OpenFileForReading(filename string) (*os.File, error) {
	return os.Open(filename)
}

// OpenFileForWriting opens a file for writing
func OpenFileForWriting(filename string) (*os.File, error) {
	return os.Create(filename)
}

// GetOutputFilename generates an output filename based on the input
func GetOutputFilename(inputFilename, defaultExtension string) string {
	if inputFilename == "-" {
		return "out" + defaultExtension
	}

	ext := filepath.Ext(inputFilename)
	return strings.TrimSuffix(inputFilename, ext) + defaultExtension
}

// FileExists checks if a file exists
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// StandardError reports an error and exits
func StandardError(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
