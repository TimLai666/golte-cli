package create

import (
	"embed"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

func CreateProject(projectName string, templates embed.FS, inCurrentDir bool, isSveltigo bool, bunPath string) {
	if bunPath == "" {
		bunPath = "bun"
	}
	var projectPath string
	if inCurrentDir {
		projectPath = "."
	} else {
		projectPath = filepath.Join("./", projectName)
		if _, err := os.Stat(projectPath); !os.IsNotExist(err) {
			log.Fatalf("Project '%s' already exists", projectName)
		}

		err := os.MkdirAll(projectPath, 0755)
		if err != nil {
			log.Fatalf("Failed to create project directory: %v", err)
		}
	}

	err := copyTemplateFiles(projectPath, "templates", templates)
	if err != nil {
		log.Fatalf("Failed to copy template files: %v", err)
	}

	// put main.go content
	mainContent := strings.Replace(mainContentTemplate, "{{projectName}}", projectName, -1)
	err = os.WriteFile(filepath.Join(projectPath, "main.go"), []byte(mainContent), 0644)
	if err != nil {
		log.Fatalf("Failed to write main.go file: %v", err)
	}

	// make router directory
	err = os.MkdirAll(filepath.Join(projectPath, "router"), 0755)
	if err != nil {
		log.Fatalf("Failed to create router directory: %v", err)
	}

	// put router.go content
	var ginContent string
	if isSveltigo {
		ginContent = strings.Replace(ginContentTemplate_sveltigo, "{{projectName}}", projectName, -1)
	} else {
		ginContent = strings.Replace(ginContentTemplate, "{{projectName}}", projectName, -1)
	}
	err = os.WriteFile(filepath.Join(projectPath, "router", "router.go"), []byte(ginContent), 0644)
	if err != nil {
		log.Fatalf("Failed to write router.go file: %v", err)
	}

	// put defineRoutes.go content
	var defineRoutesContentStr string
	if isSveltigo {
		defineRoutesContentStr = strings.Replace(defineRoutesSveltigoContent, "{{projectName}}", projectName, -1)
	} else {
		defineRoutesContentStr = strings.Replace(defineRoutesContent, "{{projectName}}", projectName, -1)
	}
	err = os.WriteFile(filepath.Join(projectPath, "router", "defineRoutes.go"), []byte(defineRoutesContentStr), 0644)
	if err != nil {
		log.Fatalf("Failed to write router.go file: %v", err)
	}

	// put package.json content
	packageJsonContent := strings.Replace(packageJsonContentTemplate, "{{projectName}}", projectName, -1)
	err = os.WriteFile(filepath.Join(projectPath, "package.json"), []byte(packageJsonContent), 0644)
	if err != nil {
		log.Fatalf("Failed to write package.json file: %v", err)
	}

	// Initialize Go module
	cmd := exec.Command("go", "mod", "init", projectName)
	cmd.Dir = projectPath
	if output, err := cmd.CombinedOutput(); err != nil {
		log.Fatalf("Failed to initialize Go module: %v\n%s", err, output)
	}

	// Run bun init
	bunCmd := exec.Command(bunPath, "init", "-y")
	bunCmd.Dir = projectPath
	if output, err := bunCmd.CombinedOutput(); err != nil {
		log.Fatalf("Failed to initialize bun: %v\n%s", err, output)
	}

	// Get Gin package
	getCmd := exec.Command("go", "get", "-u", "github.com/gin-gonic/gin")
	getCmd.Dir = projectPath
	if output, err := getCmd.CombinedOutput(); err != nil {
		log.Fatalf("Failed to get Gin package: %v\n%s", err, output)
	}

	// Get Golte package
	getCmd = exec.Command("go", "get", "-u", "github.com/nichady/golte")
	getCmd.Dir = projectPath
	if output, err := getCmd.CombinedOutput(); err != nil {
		log.Fatalf("Failed to get Golte package: %v\n%s", err, output)
	}

	// Install bun package
	bunInstallCmd := exec.Command(bunPath, "install", "golte@latest")
	bunInstallCmd.Dir = projectPath
	if output, err := bunInstallCmd.CombinedOutput(); err != nil {
		log.Fatalf("Failed to install bun package: %v\n%s", err, output)
	}

	// bunInstallCmd = exec.Command("bun", "install", "svelte@latest")
	// bunInstallCmd.Dir = projectPath
	// if output, err := bunInstallCmd.CombinedOutput(); err != nil {
	// 	log.Fatalf("Failed to install bun package: %v\n%s", err, output)
	// }
}

func copyTemplateFiles(destPath, templatePath string, templates embed.FS) error {
	entries, err := templates.ReadDir(templatePath)
	if err != nil {
		return fmt.Errorf("failed to read template directory: %v", err)
	}

	for _, entry := range entries {
		sourcePath := path.Join(templatePath, entry.Name())
		destFilePath := filepath.Join(destPath, entry.Name())

		if entry.IsDir() {
			err = os.MkdirAll(destFilePath, 0755)
			if err != nil {
				return fmt.Errorf("failed to create directory %s: %v", destFilePath, err)
			}
			err = copyTemplateFiles(destFilePath, sourcePath, templates)
			if err != nil {
				return err
			}
		} else {
			data, err := templates.ReadFile(sourcePath)
			if err != nil {
				return fmt.Errorf("failed to read template file %s: %v", sourcePath, err)
			}

			err = os.WriteFile(destFilePath, data, 0644)
			if err != nil {
				return fmt.Errorf("failed to write file %s: %v", destFilePath, err)
			}
		}
	}

	return nil
}
