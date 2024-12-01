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

	"github.com/fsnotify/fsnotify"
)

type watchPaths struct {
	configPath string
}

func WatchAndRebuild(projectPath, projectName string, startApp func(projectPath, projectName string, isSveltigo bool) *exec.Cmd, isSveltigo bool) {
	paths := &watchPaths{
		configPath: filepath.Join(projectPath, "golte.config.ts"),
	}

	var currentCmd *exec.Cmd
	processChannel := make(chan *exec.Cmd, 1)
	isRebuilding := atomic.Bool{}

	setupWatchers := func(watcher *fsnotify.Watcher) error {
		for _, watchPath := range watcher.WatchList() {
			watcher.Remove(watchPath)
		}

		dirsNotToWatch := []string{"node_modules", "dist", ".git", "build"}
		var dirsToWatch []string

		err := filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

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
			return fmt.Errorf("error walking directory tree: %v", err)
		}

		for _, dir := range dirsToWatch {
			err = watcher.Add(dir)
			if err != nil {
				log.Printf("Error adding watcher for %s: %v", dir, err)
			} else {
				fmt.Printf("Watching directory: %s\n", dir)
			}
		}

		return watcher.Add(paths.configPath)
	}

	startAndMonitor := func() bool {
		cmd := startApp(projectPath, projectName, isSveltigo)
		if cmd == nil {
			log.Println("Failed to start app, waiting for next file change...")
			return false
		}

		currentCmd = cmd
		processChannel <- cmd
		return true
	}

	startAndMonitor()

	fmt.Println("Running the project, and watching for changes...")

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("Failed to create watcher: %v", err)
	}
	defer watcher.Close()

	if err := setupWatchers(watcher); err != nil {
		log.Printf("Initial watcher setup failed: %v", err)
	}

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

	shouldIgnorePath := func(path string) bool {
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

					select {
					case cmd := <-processChannel:
						currentCmd = cmd
						if currentCmd != nil && currentCmd.Process != nil {
							_ = currentCmd.Process.Kill()
							_ = currentCmd.Wait()
						}
					default:
					}

					if err := setupWatchers(watcher); err != nil {
						log.Printf("Failed to reset watchers: %v", err)
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
