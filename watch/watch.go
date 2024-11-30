package watch

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"sync/atomic"

	"github.com/TimLai666/golte-cli/build"
	"github.com/fsnotify/fsnotify"
)

type watchPaths struct {
	configPath string
}

func WatchAndRebuild(projectPath, projectName string, startApp func(projectPath, projectName string) *exec.Cmd) {
	paths := &watchPaths{
		configPath: filepath.Join(projectPath, "golte.config.ts"),
	}

	// 使用指標預先分配一個 cmd
	var currentCmd *exec.Cmd
	processChannel := make(chan *exec.Cmd, 1)
	isRebuilding := atomic.Bool{}

	startAndMonitor := func() bool {
		if !build.BuildProject(projectPath, projectName) {
			log.Println("Build failed, waiting for next file change...")
			return false
		}

		cmd := startApp(projectPath, projectName)
		if cmd == nil {
			log.Println("Failed to start app, waiting for next file change...")
			return false
		}

		// 更新當前進程指標
		currentCmd = cmd
		processChannel <- cmd
		return true
	}

	// 首次啟動
	startAndMonitor()

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

	// 使用預計算的路徑
	watcher.Add(paths.configPath)

	// 預編譯正則表達式和預先計算的映射來加速檢查
	ignorePaths := map[string]bool{
		"node_modules": true,
		"dist":         true,
		".git":         true,
		"build":        true,
		".DS_Store":    true,
	}

	validExts := map[string]bool{
		".go":     true,
		".svelte": true,
		".css":    true,
		".html":   true,
		".ts":     true,
		".js":     true,
	}

	ignoreExts := map[string]bool{
		".tmp":  true,
		".temp": true,
		"~":     true,
	}

	// 優化 shouldIgnorePath 函數
	shouldIgnorePath := func(path string) bool {
		// 直接檢查完整路徑名稱
		for ignorePath := range ignorePaths {
			if strings.Contains(path, ignorePath) {
				return true
			}
		}

		ext := filepath.Ext(path)
		return ignoreExts[ext]
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok || shouldIgnorePath(event.Name) {
				continue
			}

			if event.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Remove|fsnotify.Rename) != 0 {
				// 使用 map 加速副檔名檢查
				ext := filepath.Ext(event.Name)
				if !validExts[ext] && ext != "" {
					continue
				}

				if !isRebuilding.CompareAndSwap(false, true) {
					continue
				}

				go func(eventName string) {
					defer func() {
						isRebuilding.Store(false)
						if r := recover(); r != nil {
							log.Printf("Failed to handle file change: %v", r)
						}
					}()

					fmt.Printf("\nFile changed: %s\nRebuilding project...\n", eventName)

					// 停止當前進程
					select {
					case cmd := <-processChannel:
						currentCmd = cmd
						if currentCmd != nil && currentCmd.Process != nil {
							_ = currentCmd.Process.Kill()
							_ = currentCmd.Wait()
						}
					default:
					}

					startAndMonitor()
				}(event.Name)
			}

		case err, ok := <-watcher.Errors:
			if !ok {
				continue
			}
			log.Printf("Watcher error: %v, continuing...", err)
		}
	}
}
