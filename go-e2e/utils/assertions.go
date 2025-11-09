package utils

import (
	"os"
	"testing"
)

func TestDirectoryExists(t *testing.T, directory string) {
	t.Helper()

	info, err := os.Stat(directory)
	if err != nil {
		if os.IsNotExist(err) {
			t.Fatalf("directory does not exist: %s", directory)
		} else {
			t.Fatalf("failed to stat directory: %v", err)
		}
	}

	if !info.IsDir() {
		t.Fatalf("%s exists but is not a directory", directory)
	}
}
