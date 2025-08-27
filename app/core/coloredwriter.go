package core

import (
	"io"

	"github.com/fatih/color"
)

type ColorWriter struct {
	Target io.Writer
	Color  string
	Prefix string
}

func (cw *ColorWriter) Write(p []byte) (int, error) {
	purpleFatih := color.New(color.FgMagenta).SprintFunc()
	coloredOutput := purpleFatih(string(p))

	return cw.Target.Write([]byte(coloredOutput))
}
