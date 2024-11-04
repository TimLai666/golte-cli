package main

import (
	"embed"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/TimLai666/golte-cli/build"
	"github.com/TimLai666/golte-cli/create"
)

//go:embed templates/*
var templates embed.FS

func main() {
	var rootCmd = &cobra.Command{
		Use:   "golte-cli",
		Short: "CLI tool for Golte projects",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Welcome to golte-cli! Use 'golte-cli new <project-name>' to create a new project.")
		},
	}

	rootCmd.AddCommand(newCmd)
	rootCmd.AddCommand(buildCmd)
	rootCmd.AddCommand(runCmd)
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
		fmt.Println("Creating project, please wait...")
		create.CreateProject(projectName, templates)
		build.BuildProject(projectName, projectName)
		fmt.Printf("Project '%s' created successfully!\n", projectName)
	},
}

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build the project",
	Run: func(cmd *cobra.Command, args []string) {
		projectPath, err := os.Getwd()
		if err != nil {
			log.Fatalf("Failed to get current directory: %v", err)
		}
		projectName := filepath.Base(projectPath)
		fmt.Println("Building the project...")
		build.BuildProject(projectPath, projectName)
	},
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Build and run the project",
	Run: func(cmd *cobra.Command, args []string) {
		projectPath, err := os.Getwd()
		if err != nil {
			log.Fatalf("Failed to get current directory: %v", err)
		}
		projectName := filepath.Base(projectPath)
		fmt.Println("Building the project...")
		build.BuildProject(projectPath, projectName)
		fmt.Println("Running the project...")

		// 創建一個新的命令
		command := exec.Command(filepath.Join("dist", projectName))

		// 將命令的標準輸出和標準錯誤直接連接到當前程序
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr

		// 執行命令
		err = command.Run()
		if err != nil {
			log.Fatalf("Failed to run project: %v", err)
		}
	},
}
