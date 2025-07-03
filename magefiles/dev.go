//go:build mage

package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/magefile/mage/mg"
)

// ─── Development ──────────────────────────────────────────────────────────────

// Dev starts development environment
func Dev() error {
	mg.Deps(Gen)

	if err := ensureTool(toolAir); err != nil {
		fmt.Println("⚠️ Air not found, using basic mode")
		return devFallback()
	}

	fmt.Println("🚀 Starting dev environment...")

	tailwindCmd, err := startTailwind()
	if err != nil {
		return fmt.Errorf("tailwind failed: %w", err)
	}
	defer cleanupTailwind(tailwindCmd)

	airCmd := setupAir()
	fmt.Println("📍 http://localhost:" + devPort)

	return runWithShutdown(airCmd, func() {
		cleanupTailwind(tailwindCmd)
	})
}

// Watch monitors templ files
func Watch() error {
	if err := ensureTool(toolTempl); err != nil {
		return err
	}

	fmt.Println("👁️ Watching...")
	return runCmd(toolTempl, "generate", "--watch")
}

// DevClean manually cleans up development artifacts
func DevClean() error {
	fmt.Println("🧹 Cleaning development environment...")

	// Kill any running development processes first
	killDevProcesses()

	// Wait a moment for processes to fully terminate
	time.Sleep(2 * time.Second)

	// Now clean up files
	cleanupTailwind(nil) // Pass nil since we're not killing any process
	return nil
}

// killDevProcesses attempts to kill common development processes
func killDevProcesses() {
	processes := []string{"air", "tailwindcss", "npx"}

	for _, proc := range processes {
		fmt.Printf("🔪 Attempting to kill %s processes...\n", proc)
		// On Windows, use taskkill
		_ = runCmdQuiet("taskkill", "/F", "/IM", proc+".exe")
		// On Unix-like systems, this would fail silently which is fine
		_ = runCmdQuiet("pkill", "-f", proc)
	}
}

// ─── Internal Functions ───────────────────────────────────────────────────────

func devFallback() error {
	mg.Deps(Gen)

	fmt.Println("🔄 Basic dev mode...")
	cmd := exec.Command("go", "run", "cmd/server/main.go")
	cmd.Env = append(os.Environ(), devEnv, "PORT="+devPort)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}

func startTailwind() (*exec.Cmd, error) {
	fmt.Println("📦 Installing Tailwind...")
	if err := runCmd("npm", "install", "-D", "tailwindcss@latest", "@tailwindcss/cli@latest"); err != nil {
		return nil, err
	}

	cmd := exec.Command("npx", "@tailwindcss/cli",
		"-i", "./internal/app.css",
		"-o", "./static/index.css",
		"--minify", "--watch")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	fmt.Printf("✅ Tailwind started (pid: %d)\n", cmd.Process.Pid)
	return cmd, nil
}

func setupAir() *exec.Cmd {
	cmd := exec.Command(toolAir, "-c", ".air.toml")
	cmd.Env = append(os.Environ(), devEnv, "PORT="+devPort)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd
}

func runWithShutdown(cmd *exec.Cmd, cleanup func()) error {
	sigs := make(chan os.Signal, 1)
	done := make(chan error, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		done <- cmd.Run()
	}()

	select {
	case sig := <-sigs:
		fmt.Printf("\nShutting down (%s)...\n", sig)
		// Kill the Air process first
		if cmd.Process != nil {
			_ = cmd.Process.Kill()
		}
		cleanup()
		return nil
	case err := <-done:
		cleanup()
		return err
	}
}

// list of files and folders to delete after dev
var packages = []string{
	// NPM/Node.js artifacts
	"node_modules",
	"package.json",
	"package-lock.json",
	".npm",
	"npm-debug.log*",
	"yarn-error.log*",

	// Development server artifacts
	"cmd/server/tmp",
	"tmp/",
	"build-errors.log",

	// Generated template files (keep source .templ files, remove generated _templ.go)
	"views/*_templ.go",

	// Build artifacts
	"bin/",
	"dist/",

	// Development environment files (optional - uncomment if you want to clean these)
	// ".env",

	// Go workspace files that might be created during dev
	"go.work",
	"go.work.sum",

	// IDE/Editor artifacts
	".vscode/settings.json",

	// OS-specific artifacts
	".DS_Store",
	"Thumbs.db",
}

func cleanupTailwind(cmd *exec.Cmd) {
	if cmd != nil && cmd.Process != nil {
		fmt.Printf("🔪 Killing Tailwind process (pid: %d)...\n", cmd.Process.Pid)
		_ = cmd.Process.Kill()
		_ = cmd.Wait()

		// Give Windows time to release native file handles
		time.Sleep(3 * time.Second)
	}
	fmt.Println("🧹 Cleaning up...")

	for _, pattern := range packages {
		// Handle glob patterns (like views/*_templ.go)
		if strings.Contains(pattern, "*") {
			cleanupGlobPattern(pattern)
		} else {
			// Handle regular files and directories
			cleanupPath(pattern)
		}
	}
}

// cleanupPath handles cleanup of individual files/directories
func cleanupPath(path string) {
	if _, err := os.Stat(path); err == nil {
		fmt.Printf("📁 Removing %s...\n", path)
		if err := removeWithRetry(path, 5); err != nil {
			fmt.Printf("⚠️ Failed to remove %s: %v\n", path, err)
		} else {
			fmt.Printf("✅ %s removed\n", path)
		}
	} else {
		fmt.Printf("ℹ️ %s not found\n", path)
	}
}

// cleanupGlobPattern handles cleanup using glob patterns
func cleanupGlobPattern(pattern string) {
	matches, err := filepath.Glob(pattern)
	if err != nil {
		fmt.Printf("⚠️ Failed to match pattern %s: %v\n", pattern, err)
		return
	}

	if len(matches) == 0 {
		fmt.Printf("ℹ️ No files match pattern %s\n", pattern)
		return
	}

	for _, match := range matches {
		fmt.Printf("📁 Removing %s (matched %s)...\n", match, pattern)
		if err := removeWithRetry(match, 5); err != nil {
			fmt.Printf("⚠️ Failed to remove %s: %v\n", match, err)
		} else {
			fmt.Printf("✅ %s removed\n", match)
		}
	}
}

// removeWithRetry attempts to remove a file/folder with retries for Windows file locking issues
func removeWithRetry(path string, maxRetries int) error {
	// Safety check - don't remove critical system paths
	if isCriticalPath(path) {
		fmt.Printf("⚠️ Skipping critical path: %s\n", path)
		return nil
	}

	var lastErr error
	for i := 0; i < maxRetries; i++ {
		if err := os.RemoveAll(path); err == nil {
			return nil
		} else {
			lastErr = err
			if i < maxRetries-1 {
				fmt.Printf("🔄 Retry %d/%d for %s...\n", i+1, maxRetries, path)
				// Wait longer between retries, especially for Windows file locks
				time.Sleep(time.Second * time.Duration(i+2))
			}
		}
	}

	// If we still can't remove it, try a more aggressive approach on Windows
	if lastErr != nil {
		fmt.Printf("🔧 Attempting forceful cleanup for %s...\n", path)
		if err := forceRemoveWindows(path); err == nil {
			return nil
		}
	}

	return lastErr
}

// isCriticalPath checks if a path is critical and should not be removed
func isCriticalPath(path string) bool {
	criticalPaths := []string{
		".", "..", "/", "\\", "C:", "C:\\",
		"go.mod", "go.sum", "main.go", "magefiles",
		".git", ".github", "README.md", "LICENSE",
	}

	for _, critical := range criticalPaths {
		if path == critical {
			return true
		}
	}

	return false
}

// forceRemoveWindows attempts to force remove files on Windows using system commands
func forceRemoveWindows(path string) error {
	// Try using Windows rmdir command with force flag
	cmd := exec.Command("cmd", "/C", "rmdir", "/S", "/Q", path)
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Run()
}
