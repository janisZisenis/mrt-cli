package log

import (
	"fmt"
	"github.com/fatih/color"
)

func Warning(message string) {
	color.Yellow(message)
}

func Error(message string) {
	color.Red(message)
}

func Success(message string) {
	color.Green(message)
}

func Info(message string) {
	fmt.Println(message)
}
