//go:build mage

package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
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
	"node_modules",
	"package.json",
	"package-lock.json",
	"cmd/server/tmp",
	"tmp/",
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

	for _, fileOrFolder := range packages {
		if _, err := os.Stat(fileOrFolder); err == nil {
			fmt.Printf("📁 Removing %s...\n", fileOrFolder)
			if err := removeWithRetry(fileOrFolder, 5); err != nil {
				fmt.Printf("⚠️ Failed to remove %s: %v\n", fileOrFolder, err)
			} else {
				fmt.Printf("✅ %s removed\n", fileOrFolder)
			}
		} else {
			fmt.Printf("ℹ️ %s not found\n", fileOrFolder)
		}
	}
}

// removeWithRetry attempts to remove a file/folder with retries for Windows file locking issues
func removeWithRetry(path string, maxRetries int) error {
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

// forceRemoveWindows attempts to force remove files on Windows using system commands
func forceRemoveWindows(path string) error {
	// Try using Windows rmdir command with force flag
	cmd := exec.Command("cmd", "/C", "rmdir", "/S", "/Q", path)
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Run()
}
