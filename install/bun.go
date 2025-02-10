package install

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func InstallBun() (bunPath string, err error) {
	// 檢查是否已安裝
	if path := getBunPath(); path != "" {
		return path, nil
	}

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("powershell", "-c", "irm bun.sh/install.ps1 | iex")
		err = cmd.Run()
		if err != nil {
			return "", fmt.Errorf("failed to install bun: %v", err)
		}
		return getBunPath(), nil
	case "linux", "darwin":
		path := findBunInUnix()
		if path != "" {
			fmt.Println("Bun already installed at", path)
			return path, nil
		}
		cmd = exec.Command("bash", "-c", "curl -fsSL https://bun.sh/install | bash")
		err = cmd.Run()
		if err != nil {
			return "", fmt.Errorf("failed to install bun: %v", err)
		}
		return findBunInUnix(), nil
	default:
		return "", fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
}

func getBunPath() string {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("where", "bun.exe")
	} else {
		cmd = exec.Command("which", "bun")
	}
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}

func findBunInUnix() (path string) {
	homeDir := os.Getenv("HOME")
	if homeDir == "" {
		return ""
	}
	cmd := exec.Command("find", homeDir, "-name", "bun", "-type", "f")
	b, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(b))
}
