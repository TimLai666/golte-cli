package main

import (
	"embed"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

//go:embed templates/*
var templates embed.FS

func main() {
	var rootCmd = &cobra.Command{
		Use:   "golte-cli",
		Short: "CLI tool for generating Golte sample projects",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Welcome to golte-cli! Use 'golte-cli new <project-name>' to create a new project.")
		},
	}

	rootCmd.AddCommand(newCmd)
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Error executing command: %v", err)
	}
}

var newCmd = &cobra.Command{
	Use:   "new <project-name>",
	Short: "Create a new Golte sample project",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectName := args[0]
		createProject(projectName)
	},
}

func createProject(projectName string) {
	projectPath := filepath.Join("./", projectName)
	if _, err := os.Stat(projectPath); !os.IsNotExist(err) {
		log.Fatalf("Project '%s' already exists", projectName)
	}

	err := os.MkdirAll(projectPath, 0755)
	if err != nil {
		log.Fatalf("Failed to create project directory: %v", err)
	}

	err = copyTemplateFiles(projectPath, "templates")
	if err != nil {
		log.Fatalf("Failed to copy template files: %v", err)
	}

	// Initialize Go module
	cmd := exec.Command("go", "mod", "init", projectName)
	cmd.Dir = projectPath
	if output, err := cmd.CombinedOutput(); err != nil {
		log.Fatalf("Failed to initialize Go module: %v\n%s", err, output)
	}

	// Run npm init
	npmCmd := exec.Command("npm", "init", "-y")
	npmCmd.Dir = projectPath
	if output, err := npmCmd.CombinedOutput(); err != nil {
		log.Fatalf("Failed to initialize npm: %v\n%s", err, output)
	}

	// Get Golte package
	getCmd := exec.Command("go", "get", "-u", "github.com/gin-gonic/gin")
	getCmd.Dir = projectPath
	if output, err := getCmd.CombinedOutput(); err != nil {
		log.Fatalf("Failed to get Gin package: %v\n%s", err, output)
	}

	getCmd = exec.Command("go", "get", "-u", "github.com/nichady/golte")
	getCmd.Dir = projectPath
	if output, err := getCmd.CombinedOutput(); err != nil {
		log.Fatalf("Failed to get Golte package: %v\n%s", err, output)
	}

	// Install npm package
	npmInstallCmd := exec.Command("npm", "install", "golte@latest")
	npmInstallCmd.Dir = projectPath
	if output, err := npmInstallCmd.CombinedOutput(); err != nil {
		log.Fatalf("Failed to install npm package: %v\n%s", err, output)
	}

	fmt.Printf("Project '%s' created successfully!\n", projectName)
}

func copyTemplateFiles(destPath, templatePath string) error {
	entries, err := templates.ReadDir(templatePath)
	if err != nil {
		return fmt.Errorf("failed to read template directory: %v", err)
	}

	for _, entry := range entries {
		sourcePath := filepath.Join(templatePath, entry.Name())
		destFilePath := filepath.Join(destPath, entry.Name())

		if entry.IsDir() {
			err = os.MkdirAll(destFilePath, 0755)
			if err != nil {
				return fmt.Errorf("failed to create directory %s: %v", destFilePath, err)
			}
			err = copyTemplateFiles(destFilePath, sourcePath)
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
