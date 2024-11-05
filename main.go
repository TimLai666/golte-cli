package main

import (
	"embed"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/fsnotify/fsnotify"

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
			fmt.Println("Welcome to golte-cli! Use `golte-cli help` to see available commands.")
		},
	}

	rootCmd.AddCommand(newCmd)
	rootCmd.AddCommand(buildCmd)
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(devCmd)
	rootCmd.HelpFunc()
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

var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "Run the project and auto rebuild when changes",
	Run: func(cmd *cobra.Command, args []string) {
		projectPath, err := os.Getwd()
		if err != nil {
			log.Fatalf("Failed to get current directory: %v", err)
		}
		projectName := filepath.Base(projectPath)

		// 創建一個通道用於進程管理
		processChannel := make(chan *exec.Cmd, 1)

		// 定義啟動應用程序的函數
		startApp := func() *exec.Cmd {
			build.BuildProject(projectPath, projectName)
			cmd := exec.Command(filepath.Join("dist", projectName))
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Start(); err != nil {
				log.Printf("Failed to start project: %v", err)
				return nil
			}
			return cmd
		}

		// 首次啟動應用程序
		if cmd := startApp(); cmd != nil {
			processChannel <- cmd
		}

		fmt.Println("Running the project, and watching for changes...")

		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			log.Fatalf("Failed to create watcher: %v", err)
		}
		defer watcher.Close()

		// 監控目錄設置
		dirsNotToWatch := []string{"node_modules", "dist", ".git", "build"}
		var dirsToWatch []string

		err = filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// 檢查是否應該跳過此目錄
			if info.IsDir() {
				baseName := info.Name()
				if slices.Contains(dirsNotToWatch, baseName) {
					fmt.Printf("Skipping directory: %s\n", path)
					return filepath.SkipDir
				}
				dirsToWatch = append(dirsToWatch, path)
			}
			return nil
		})

		if err != nil {
			log.Printf("Error walking directory tree: %v", err)
		}

		// 為每個要監控的目錄添加監控
		for _, dir := range dirsToWatch {
			err = watcher.Add(dir)
			if err != nil {
				log.Printf("Error adding watcher for %s: %v", dir, err)
			} else {
				fmt.Printf("Watching directory: %s\n", dir)
			}
		}

		// 添加根目錄的特定文件
		watcher.Add(filepath.Join(projectPath, "golte.config.ts"))

		debounceTimer := time.NewTimer(0)
		debounceTimer.Stop()

		// 定義檢查路徑是否應該被忽略的函數
		shouldIgnorePath := func(path string) bool {
			ignorePaths := []string{
				"node_modules",
				"dist",
				".git",
				"build",
				"temp-",     // 臨時文件
				".DS_Store", // macOS 系統文件
			}

			// 檢查路徑是否包含任何需要忽略的部分
			for _, ignore := range ignorePaths {
				if strings.Contains(path, ignore) {
					return true
				}
			}

			// 檢查文件擴展名
			ext := filepath.Ext(path)
			ignoreExts := []string{".tmp", ".temp", "~"}
			for _, ignoreExt := range ignoreExts {
				if strings.HasSuffix(ext, ignoreExt) {
					return true
				}
			}

			return false
		}

		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				// 使用改進的忽略檢查
				if shouldIgnorePath(event.Name) {
					continue
				}

				if event.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Remove|fsnotify.Rename) != 0 {
					// 只處理特定文件類型
					ext := filepath.Ext(event.Name)
					validExts := []string{".go", ".svelte", ".css", ".html", ".ts", ".js"}
					isValidExt := false
					for _, validExt := range validExts {
						if ext == validExt {
							isValidExt = true
							break
						}
					}

					if !isValidExt && ext != "" { // 允許沒有擴展名的文件（如配置文件）
						continue
					}

					debounceTimer.Reset(500 * time.Millisecond)
					go func() {
						<-debounceTimer.C
						fmt.Printf("\nFile changed: %s\n", event.Name)
						fmt.Println("Rebuilding project...")

						// 停止當前運行的進程
						if currentCmd := <-processChannel; currentCmd != nil && currentCmd.Process != nil {
							if err := currentCmd.Process.Kill(); err != nil {
								log.Printf("Error killing process: %v", err)
							}
							currentCmd.Wait()
						}

						// 啟動新的應用程序
						if newCmd := startApp(); newCmd != nil {
							processChannel <- newCmd
						}
					}()
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("Watcher error:", err)
			}
		}
	},
}
