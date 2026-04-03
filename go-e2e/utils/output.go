package utils

import (
	"regexp"
	"strings"
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

func (o *Output) AssertLineEquals(t *testing.T, index int, expectedText string) {
	t.Helper()

	require.Less(t, index, len(o.lines), "line index %d is out of bounds, have %d lines", index, len(o.lines))
	assert.Equal(t, expectedText, o.lines[index], "line %d does not contain expected text", index)
}

func (o *Output) AssertLineMatchesRegex(t *testing.T, index int, pattern string) {
	t.Helper()

	require.Less(t, index, len(o.lines), "line index %d is out of bounds, have %d lines", index, len(o.lines))
	regex, err := regexp.Compile(pattern)
	require.NoError(t, err, "invalid regex pattern: %s", pattern)
	assert.True(t, regex.MatchString(o.lines[index]), "line %d does not match pattern %s\ngot: %s", index, pattern, o.lines[index])
}

func (o *Output) AssertHasLine(t *testing.T, line string) {
	t.Helper()

	assert.Contains(t, o.lines, line, "output does not have line: %s", line)
}

func (o *Output) AssertNextLineAfterLineContaining(t *testing.T, containing string, expectedNextLine string) {
	t.Helper()

	for i, line := range o.lines {
		if strings.Contains(line, containing) {
			require.Less(t, i+1, len(o.lines), "no line after line containing: %s", containing)
			assert.Equal(t, expectedNextLine, o.lines[i+1], "line after line containing %q does not equal expected text", containing)
			return
		}
	}
	assert.Fail(t, "output does not have a line containing: "+containing)
}
