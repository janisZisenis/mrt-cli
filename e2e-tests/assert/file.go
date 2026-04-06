package assert

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func FileExists(t *testing.T, path string) {
	t.Helper()

	_, err := os.Stat(path)
	assert.False(t, os.IsNotExist(err), "file %s does not exist but should", path)
}

func FileDoesNotExist(t *testing.T, path string) {
	t.Helper()

	_, err := os.Stat(path)
	assert.True(t, os.IsNotExist(err), "file %s exists but should not", path)
}

func FileHasContent(t *testing.T, path string, expectedContent string) {
	t.Helper()

	content, err := os.ReadFile(path)
	require.NoError(t, err, "failed to read file: %s", path)
	assert.Equal(t, expectedContent, string(content))
}
