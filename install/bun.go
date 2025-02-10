package install

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func InstallBun() (bunPath string, err error) {
	if hasBunInstalled() {
		return "", nil
	}
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("powershell", "-c", "irm bun.sh/install.ps1 | iex")
	case "linux", "darwin":
		path := findBunInUnix()
		if path != "" {
			fmt.Println("Bun already installed at", path)
			return path, nil
		}
		cmd = exec.Command("bash", "-c", "curl -fsSL https://bun.sh/install | bash")
		err := cmd.Run()
		if err != nil {
			return "", err
		}
		return findBunInUnix(), nil

	}
	err = cmd.Run()
	return "", err
}

func hasBunInstalled() bool {
	cmd := exec.Command("bun", "--version")
	err := cmd.Run()
	return err == nil
}

func findBunInUnix() (path string) {
	homeDir := os.Getenv("HOME")
	cmd := exec.Command("find", homeDir, "-name", "bun", "-type", "f")
	b, err := cmd.Output()
	if err != nil {
		return ""
	}
	// 移除結尾的換行符號
	return strings.TrimSpace(string(b))
}
