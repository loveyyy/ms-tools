package utils

import (
	"os"
	"os/exec"
	"path"
	"path/filepath"
)

func GetRootPath() string {
	filePath, _ := exec.LookPath(os.Args[0])
	absFilePath, _ := filepath.Abs(filePath)
	return path.Dir(absFilePath)
}
