package install

import (
	"os/exec"
	"runtime"
)

func InstallBun() {
	if hasBunInstalled() {
		return
	}
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("powershell", "-c", "irm bun.sh/install.ps1 | iex")
	case "linux", "darwin":
		cmd = exec.Command("bash", "-c", "curl -fsSL https://bun.sh/install | bash")
	}
	cmd.Run()
}

func hasBunInstalled() bool {
	cmd := exec.Command("bun", "--version")
	err := cmd.Run()
	return err == nil
}
