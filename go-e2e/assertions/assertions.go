package assertions

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDirectoryExists(t *testing.T, directory string) {
	t.Helper()

	info, err := os.Stat(directory)
	assert.NoError(t, err, "failed to stat directory: %s", directory)
	if err != nil {
		return
	}

	assert.True(t, info.IsDir(), "%s exists but is not a directory", directory)
}
