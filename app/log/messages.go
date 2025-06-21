package log

import (
	"fmt"
	"github.com/fatih/color"
)

func Warning(format string, args ...interface{}) {
	color.Yellow(formatMessage(format, args))
}

func Error(format string, args ...interface{}) {
	color.Red(formatMessage(format, args))
}

func Success(format string, args ...interface{}) {
	color.Green(formatMessage(format, args))
}

func Info(format string, args ...interface{}) {
	fmt.Println(formatMessage(format, args))
}

func formatMessage(format string, args []interface{}) string {
	return fmt.Sprintf(format, args...)
}
