package output_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"mrt-cli/go-e2e/output"
)

// --- ExactLine ---

func Test_exactLine_Matches_IdenticalString(t *testing.T) {
	assert.True(t, output.ExactLine{"hello"}.Matches("hello"))
}

func Test_exactLine_DoesNotMatch_WhenLineContainsText(t *testing.T) {
	assert.False(t, output.ExactLine{"hello"}.Matches("say hello world"))
}

func Test_exactLine_DoesNotMatch_WhenLineIsDifferent(t *testing.T) {
	assert.False(t, output.ExactLine{"hello"}.Matches("world"))
}

// --- ContainsLine ---

func Test_containsLine_Matches_IdenticalString(t *testing.T) {
	assert.True(t, output.ContainsLine{"hello"}.Matches("hello"))
}

func Test_containsLine_Matches_WhenLineContainsText(t *testing.T) {
	assert.True(t, output.ContainsLine{"hello"}.Matches("say hello world"))
}

func Test_containsLine_DoesNotMatch_WhenLineDoesNotContainText(t *testing.T) {
	assert.False(t, output.ContainsLine{"hello"}.Matches("goodbye"))
}
