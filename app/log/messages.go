package log

import (
	"fmt"
	"github.com/fatih/color"
)

func Warningf(format string, args ...interface{}) {
	color.Yellow(formatMessage(format, args))
}

func Errorf(format string, args ...interface{}) {
	color.Red(formatMessage(format, args))
}

func Successf(format string, args ...interface{}) {
	color.Green(formatMessage(format, args))
}

func Infof(format string, args ...interface{}) {
	fmt.Println(formatMessage(format, args))
}

func formatMessage(format string, args []interface{}) string {
	return fmt.Sprintf(format, args...)
}
