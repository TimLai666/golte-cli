package watch

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

func WatchAndRebuild(projectPath, projectName string, startApp func(projectPath, projectName string) *exec.Cmd) {
	// 創建一個通道用於進程管理
	processChannel := make(chan *exec.Cmd, 1)

	// 創建一個通道用於監控進程狀態
	processExitChannel := make(chan struct{}, 1)

	// 啟動應用程序
	startAndMonitor := func() {
		for {
			cmd := startApp(projectPath, projectName)
			if cmd == nil {
				log.Println("啟動失敗，5秒後重試...")
				time.Sleep(5 * time.Second)
				continue // 直接重試，不使用 channel
			}

			// 成功啟動
			processChannel <- cmd
			go func() {
				defer func() {
					if r := recover(); r != nil {
						log.Printf("監控進程時發生錯誤: %v", r)
						processExitChannel <- struct{}{}
					}
				}()

				cmd.Wait()
				processExitChannel <- struct{}{}
			}()
			break
		}
	}

	// 啟動監控循環
	go func() {
		for {
			startAndMonitor()
			// 等待進程退出信號
			<-processExitChannel
			// log.Println("檢測到進程退出，5秒後重新啟動...")
			time.Sleep(5 * time.Second)
		}
	}()

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

	// 定檢查路徑是否應該被忽略的函數
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
				continue
			}

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

				if !isValidExt && ext != "" {
					continue
				}

				debounceTimer.Reset(500 * time.Millisecond)
				go func(eventName string) {
					defer func() {
						if r := recover(); r != nil {
							log.Printf("處理文件變更時發生錯誤: %v", r)
							processExitChannel <- struct{}{}
						}
					}()

					<-debounceTimer.C
					fmt.Printf("\nFile changed: %s\n", eventName)
					fmt.Println("Rebuilding project...")

					// 停止當前進程
					if currentCmd := <-processChannel; currentCmd != nil && currentCmd.Process != nil {
						_ = currentCmd.Process.Kill()
						_ = currentCmd.Wait()
					}

					// 重新啟動
					startAndMonitor()
				}(event.Name)
			}

		case err, ok := <-watcher.Errors:
			if !ok {
				continue
			}
			log.Printf("監視器錯誤: %v，繼續運行...", err)
		}
	}
}
