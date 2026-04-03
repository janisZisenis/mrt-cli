package utils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func AssertDirectoryExists(t *testing.T, directory string) {
	t.Helper()

	info, err := os.Stat(directory)
	require.NoError(t, err, "failed to stat directory: %s", directory)
	assert.True(t, info.IsDir(), "%s exists but is not a directory", directory)
}

func AssertDirectoryDoesNotExist(t *testing.T, directory string) {
	t.Helper()

	_, err := os.Stat(directory)
	assert.True(t, os.IsNotExist(err), "directory %s exists but should not", directory)
}
