package tests_test

import (
	"mrt-cli/e2e-tests/fixtures"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_IfTeamDirDoesNotExist_ShouldPrintErrorMessageAndExitWithError(t *testing.T) {
	f := fixtures.MakeMrtFixture(t).Parallel()
	notExisting := "/some/unknown/path"

	output, exitCode := f.MakeMrtCommand().
		RunInDirectory(notExisting).
		Setup().
		All().
		Execute()

	assert.NotEqual(t, 0, exitCode)
	output.AssertHasLine(t, "Directory \""+notExisting+"\" does not exist.")
}
