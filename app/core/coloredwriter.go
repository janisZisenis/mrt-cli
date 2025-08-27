package core

import (
	"fmt"
	"github.com/fatih/color"
	"io"
)

type ColorWriter struct {
	Target io.Writer
	Color  string
	Prefix string
}

func (cw *ColorWriter) Write(p []byte) (n int, err error) {
	purpleFatih := color.New(color.FgMagenta).SprintFunc()

	coloredOutput := fmt.Sprintf("%s", purpleFatih(string(p)))
	return cw.Target.Write([]byte(coloredOutput))
}
