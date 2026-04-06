package outputs_test

import (
	"testing"

	"mrt-cli/e2e-tests/outputs"

	"github.com/stretchr/testify/assert"
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
	o := outputs.Make([]string{"line A", "line B"})
	o.AssertInOrder(t, outputs.HasLine("line A"))
}

func Test_AssertInOrder_Passes_WhenMultipleLinesFoundInOrder(t *testing.T) {
	o := outputs.Make([]string{"line A", "line B", "line C"})
	o.AssertInOrder(t, outputs.HasLine("line A"), outputs.HasLine("line C"))
}

func Test_AssertInOrder_Passes_WithNonAdjacentLines(t *testing.T) {
	o := outputs.Make([]string{"line A", "irrelevant", "line B"})
	o.AssertInOrder(t, outputs.HasLine("line A"), outputs.HasLine("line B"))
}

func Test_AssertInOrder_Passes_WithHasLineContaining(t *testing.T) {
	o := outputs.Make([]string{"cloning repo", "clone failed: permission denied", "skipping"})
	o.AssertInOrder(t,
		outputs.HasLine("cloning repo"),
		outputs.HasLineContaining("clone failed:"),
		outputs.HasLine("skipping"),
	)
}

// --- AssertInOrder failure paths ---

func Test_AssertInOrder_Fails_WhenLinesInWrongOrder(t *testing.T) {
	o := outputs.Make([]string{"line A", "line B"})
	spy := runWithSpy(func(spy *spyT) {
		o.AssertInOrder(spy, outputs.HasLine("line B"), outputs.HasLine("line A"))
	})
	assert.True(t, spy.failNowCalled)
}

func Test_AssertInOrder_Fails_WhenLineNotFound(t *testing.T) {
	o := outputs.Make([]string{"line A", "line B"})
	spy := runWithSpy(func(spy *spyT) {
		o.AssertInOrder(spy, outputs.HasLine("line C"))
	})
	assert.True(t, spy.failNowCalled)
}

func Test_AssertInOrder_Fails_WhenContainingLineNotFound(t *testing.T) {
	o := outputs.Make([]string{"line A", "line B"})
	spy := runWithSpy(func(spy *spyT) {
		o.AssertInOrder(spy, outputs.HasLineContaining("not present"))
	})
	assert.True(t, spy.failNowCalled)
}
