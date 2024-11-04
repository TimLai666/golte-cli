package build

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"
)

func BuildProject(projectPath string, projectName string) {
	// build frontend
	cmd := exec.Command("npx", "golte")
	cmd.Dir = projectPath
	if output, err := cmd.CombinedOutput(); err != nil {
		log.Fatalf("Failed to build frontend: %v\n%s", err, output)
	}

	// build the project
	var execName string
	if runtime.GOOS == "windows" {
		execName = fmt.Sprintf("%s.exe", projectName)
	} else {
		execName = projectName
	}
	cmd = exec.Command("go", "build", "-o", execName, "main.go")
	cmd.Dir = projectPath
	if output, err := cmd.CombinedOutput(); err != nil {
		log.Fatalf("Failed to build project: %v\n%s", err, output)
	}
}
