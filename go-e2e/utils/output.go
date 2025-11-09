package utils

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/slices"
)

type Output struct {
	lines []string
}

func MakeOutput(lines []string) *Output {
	return &Output{
		lines: lines,
	}
}

func (o *Output) Reversed() *Output {
	reversed := append([]string(nil), o.lines...)
	slices.Reverse(reversed)

	return MakeOutput(reversed)
}

func (o *Output) AssertLineContains(t *testing.T, index int, expectedText string) {
	t.Helper()

	require.Less(t, index, len(o.lines), "line index %d is out of bounds, have %d lines", index, len(o.lines))
	assert.Contains(t, o.lines[index], expectedText, "line %d does not contain expected text", index)
}

func (o *Output) AssertLineMatchesRegex(t *testing.T, index int, pattern string) {
	t.Helper()

	assert.Less(t, index, len(o.lines), "line index %d is out of bounds, have %d lines", index, len(o.lines))
	if index >= len(o.lines) {
		return
	}

	regex, err := regexp.Compile(pattern)
	assert.NoError(t, err, "invalid regex pattern: %s", pattern)
	if err != nil {
		return
	}

	assert.True(t, regex.MatchString(o.lines[index]), "line %d does not match pattern %s\ngot: %s", index, pattern, o.lines[index])
}
