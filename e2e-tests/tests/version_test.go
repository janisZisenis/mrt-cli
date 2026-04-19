package tests_test

import (
	"mrt-cli/e2e-tests/fixtures"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_VersionCommand_PrintsVersionInfo(t *testing.T) {
	f := fixtures.MakeMrtFixture(t)

	output, exitCode := f.MakeMrtCommand().
		Version().
		Execute()

	assert.Equal(t, 0, exitCode)
	output.AssertLineMatchesRegex(t, 0, `^mrt - version .+, commit .+, built at .+ by .+$`)
}
