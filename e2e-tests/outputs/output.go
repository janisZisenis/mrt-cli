package outputs

import (
	"fmt"
	"regexp"
	"slices"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Output struct {
	lines []string
}

func Make(lines []string) *Output {
	return &Output{
		lines: lines,
	}
}

func (o *Output) dump() string {
	var sb strings.Builder
	sb.WriteString("captured output:\n")
	for i, line := range o.lines {
		fmt.Fprintf(&sb, "  [%d] %s\n", i, line)
	}
	return sb.String()
}

func (o *Output) Reversed() *Output {
	reversed := append([]string(nil), o.lines...)
	slices.Reverse(reversed)

	return Make(reversed)
}

func (o *Output) AssertLineEquals(t *testing.T, index int, expectedText string) {
	t.Helper()

	require.Less(
		t, index, len(o.lines),
		"line index %d is out of bounds, have %d lines\n%s", index, len(o.lines), o.dump(),
	)
	assert.Equal(t, expectedText, o.lines[index], "line %d does not match expected text\n%s", index, o.dump())
}

func (o *Output) AssertLineMatchesRegex(t *testing.T, index int, pattern string) {
	t.Helper()

	require.Less(
		t, index, len(o.lines),
		"line index %d is out of bounds, have %d lines\n%s", index, len(o.lines), o.dump(),
	)
	regex, err := regexp.Compile(pattern)
	require.NoError(t, err, "invalid regex pattern: %s", pattern)
	assert.True(
		t,
		regex.MatchString(o.lines[index]),
		"line %d does not match pattern %s\n%s",
		index,
		pattern,
		o.dump(),
	)
}

func (o *Output) AssertHasLine(t *testing.T, line string) {
	t.Helper()

	assert.Contains(t, o.lines, line, "output does not have line: %s\n%s", line, o.dump())
}

func (o *Output) AssertHasNoLineContaining(t *testing.T, text string) {
	t.Helper()

	for _, line := range o.lines {
		if strings.Contains(line, text) {
			assert.Fail(t, "output should not contain: "+text+"\n"+o.dump())
			return
		}
	}
}

type LineExpectation interface {
	matches(line string) bool
	failureMessage() string
}

func (o *Output) AssertInOrder(t require.TestingT, expectations ...LineExpectation) {
	if h, ok := t.(interface{ Helper() }); ok {
		h.Helper()
	}
	cursor := 0
	for _, exp := range expectations {
		found := false
		for i := cursor; i < len(o.lines); i++ {
			if exp.matches(o.lines[i]) {
				cursor = i + 1
				found = true
				break
			}
		}
		if !found {
			require.Fail(t, exp.failureMessage(), o.dump())
		}
	}
}
