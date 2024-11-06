package build

import (
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
	"runtime"
)

func BuildProject(projectPath string, projectName string) bool {
	// build frontend
	cmd := exec.Command("npx", "golte")
	cmd.Dir = projectPath
	if output, err := cmd.CombinedOutput(); err != nil {
		log.Printf("Failed to build frontend: %v\n%s", err, output)
		return false
	}

	// build the project
	var execName string
	if runtime.GOOS == "windows" {
		execName = fmt.Sprintf("%s.exe", projectName)
	} else {
		execName = projectName
	}
	cmd = exec.Command("go", "build", "-o", filepath.Join("dist", execName), "main.go")
	cmd.Dir = projectPath
	if output, err := cmd.CombinedOutput(); err != nil {
		log.Printf("Failed to build project: %v\n%s", err, output)
		return false
	}

	return true
}
