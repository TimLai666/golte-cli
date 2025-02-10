package install

import (
	"os/exec"
	"runtime"
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
	cmd := exec.Command("find", "$HOME", "-name", "bun", "-type", "f")
	err := cmd.Run()
	if err != nil {
		return ""
	}
	b, err := cmd.Output()
	if err != nil {
		return ""
	}
	return string(b)
}
