package utils_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"mrt-cli/go-e2e/utils"
)

// --- spy infrastructure ---

type failNowSignal struct{}

type spyT struct {
	failNowCalled bool
}

func (s *spyT) Errorf(_ string, _ ...interface{}) {}
func (s *spyT) FailNow()                          { s.failNowCalled = true; panic(failNowSignal{}) }

func runWithSpy(f func(t *spyT)) *spyT {
	spy := &spyT{}
	func() {
		defer func() {
			r := recover()
			if r == nil || r == (failNowSignal{}) {
				return
			}
			panic(r) // unexpected panic — re-panic so the test fails loudly
		}()
		f(spy)
	}()
	return spy
}

// --- AssertInOrder happy paths ---

func Test_AssertInOrder_Passes_WhenSingleLineFound(t *testing.T) {
	output := utils.MakeOutput([]string{"line A", "line B"})
	output.AssertInOrder(t, utils.HasLine("line A"))
}

func Test_AssertInOrder_Passes_WhenMultipleLinesFoundInOrder(t *testing.T) {
	output := utils.MakeOutput([]string{"line A", "line B", "line C"})
	output.AssertInOrder(t, utils.HasLine("line A"), utils.HasLine("line C"))
}

func Test_AssertInOrder_Passes_WithNonAdjacentLines(t *testing.T) {
	output := utils.MakeOutput([]string{"line A", "irrelevant", "line B"})
	output.AssertInOrder(t, utils.HasLine("line A"), utils.HasLine("line B"))
}

func Test_AssertInOrder_Passes_WithHasLineContaining(t *testing.T) {
	output := utils.MakeOutput([]string{"cloning repo", "clone failed: permission denied", "skipping"})
	output.AssertInOrder(t,
		utils.HasLine("cloning repo"),
		utils.HasLineContaining("clone failed:"),
		utils.HasLine("skipping"),
	)
}

// --- AssertInOrder failure paths ---

func Test_AssertInOrder_Fails_WhenLinesInWrongOrder(t *testing.T) {
	output := utils.MakeOutput([]string{"line A", "line B"})
	spy := runWithSpy(func(spy *spyT) {
		output.AssertInOrder(spy, utils.HasLine("line B"), utils.HasLine("line A"))
	})
	assert.True(t, spy.failNowCalled)
}

func Test_AssertInOrder_Fails_WhenLineNotFound(t *testing.T) {
	output := utils.MakeOutput([]string{"line A", "line B"})
	spy := runWithSpy(func(spy *spyT) {
		output.AssertInOrder(spy, utils.HasLine("line C"))
	})
	assert.True(t, spy.failNowCalled)
}

func Test_AssertInOrder_Fails_WhenContainingLineNotFound(t *testing.T) {
	output := utils.MakeOutput([]string{"line A", "line B"})
	spy := runWithSpy(func(spy *spyT) {
		output.AssertInOrder(spy, utils.HasLineContaining("not present"))
	})
	assert.True(t, spy.failNowCalled)
}
