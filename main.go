package main

import (
	"embed"
	"fmt"
	"log"

	"github.com/spf13/cobra"

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
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Error executing command: %v", err)
	}

	rootCmd.AddCommand(buildCmd)
}

var newCmd = &cobra.Command{
	Use:   "new <project-name>",
	Short: "Create a new Golte sample project",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectName := args[0]
		fmt.Println("Creating project, please wait...")
		create.CreateProject(projectName, templates)
	},
}

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build the project",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Building the project...")
	},
}
