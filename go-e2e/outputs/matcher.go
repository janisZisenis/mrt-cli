package outputs

import "strings"

type ExactLine struct{ Text string }

func (e ExactLine) Matches(line string) bool { return line == e.Text }
func (e ExactLine) matches(line string) bool { return e.Matches(line) }
func (e ExactLine) failureMessage() string {
	return "no line equaling \"" + e.Text + "\" found after previous matches"
}

type ContainsLine struct{ Text string }

func (c ContainsLine) Matches(line string) bool { return strings.Contains(line, c.Text) }
func (c ContainsLine) matches(line string) bool { return c.Matches(line) }
func (c ContainsLine) failureMessage() string {
	return "no line containing \"" + c.Text + "\" found after previous matches"
}

func HasLine(text string) LineExpectation           { return ExactLine{text} }
func HasLineContaining(text string) LineExpectation { return ContainsLine{text} }
