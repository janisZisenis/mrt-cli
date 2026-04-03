package outputs_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"mrt-cli/go-e2e/outputs"
)

// --- ExactLine ---

func Test_exactLine_Matches_IdenticalString(t *testing.T) {
	assert.True(t, outputs.ExactLine{"hello"}.Matches("hello"))
}

func Test_exactLine_DoesNotMatch_WhenLineContainsText(t *testing.T) {
	assert.False(t, outputs.ExactLine{"hello"}.Matches("say hello world"))
}

// --- ContainsLine ---

func Test_containsLine_Matches_IdenticalString(t *testing.T) {
	assert.True(t, outputs.ContainsLine{"hello"}.Matches("hello"))
}

func Test_containsLine_Matches_WhenLineContainsText(t *testing.T) {
	assert.True(t, outputs.ContainsLine{"hello"}.Matches("say hello world"))
}

func Test_containsLine_DoesNotMatch_WhenLineDoesNotContainText(t *testing.T) {
	assert.False(t, outputs.ContainsLine{"hello"}.Matches("goodbye"))
}
