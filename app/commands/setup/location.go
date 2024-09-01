package setup

import (
	"os"
	"path"
)

func GetExecutablePath() string {
	executable, _ := os.Executable()
	return path.Dir(executable)
}
