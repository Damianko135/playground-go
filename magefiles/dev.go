//go:build mage

package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

// ─── Development ──────────────────────────────────────────────────────────────

// Dev starts development server with hot reload and Tailwind, falls back to basic mode if Air is missing
func Dev() error {
	if err := ensureTool(toolTempl); err != nil {
		return err
	}
	if err := ensureTool(toolAir); err != nil {
		fmt.Println("⚠️ Air not found, falling back to basic dev mode")
		return devBasic()
	}

	// Install Tailwind dependencies for development
	fmt.Println("📦 Installing Tailwind CSS dependencies...")
	if err := runCmd("npm", "install", "-D", "tailwindcss@latest", "@tailwindcss/cli@latest"); err != nil {
		return fmt.Errorf("failed to install Tailwind dependencies: %w", err)
	}

	fmt.Println("🛠️ Generating templ files...")
	if err := runCmd(toolTempl, "generate"); err != nil {
		return err
	}

	tailwindCmd := exec.Command("npx", "@tailwindcss/cli", "-i", "./internal/app.css", "-o", "./static/index.css", "--minify", "--watch")
	tailwindCmd.Stdout = os.Stdout
	tailwindCmd.Stderr = os.Stderr
	tailwindCmd.Env = os.Environ()
	tailwindCmd.Dir = "." // Run from project root

	fmt.Println("📦 Running Tailwind build with npx...")
	if err := tailwindCmd.Start(); err != nil {
		return fmt.Errorf("failed to start tailwind watcher via npx: %w", err)
	}
	fmt.Printf("✅ Tailwind watcher started (pid: %d)\n", tailwindCmd.Process.Pid)

	airCmd := exec.Command(toolAir, "-c", "../../.air.toml")
	airCmd.Dir = "cmd/server"
	airCmd.Env = append(os.Environ(), "GO_ENV=development", "PORT=3000")
	airCmd.Stdout = os.Stdout
	airCmd.Stderr = os.Stderr
	airCmd.Stdin = os.Stdin

	fmt.Println("🚀 Starting development server with hot reload...")
	fmt.Println("📍 http://localhost:3000")

	sigs := make(chan os.Signal, 1)
	done := make(chan error, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		done <- airCmd.Run()
	}()

	// Helper function to cleanup processes and dependencies
	cleanup := func() {
		_ = tailwindCmd.Process.Kill()
		_ = airCmd.Process.Kill()
		fmt.Println("🧹 Cleaning up node_modules...")
		_ = os.RemoveAll("node_modules")
		_ = os.Remove("package.json")
		_ = os.Remove("package-lock.json")
	}

	select {
	case sig := <-sigs:
		fmt.Printf("\nReceived signal: %s, shutting down...\n", sig)
		cleanup()
		return nil
	case err := <-done:
		cleanup()
		return err
	}
}

// devBasic is an internal fallback function for basic development mode
func devBasic() error {
	fmt.Println("🛠️ Generating templ files...")
	if err := runCmd(toolTempl, "generate"); err != nil {
		return err
	}

	fmt.Println("🔄 Running in basic dev mode...")
	cmd := exec.Command("go", "run", ".")
	cmd.Dir = "cmd/server"
	cmd.Env = append(os.Environ(), "GO_ENV=development", "PORT=3000")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

// Watch monitors templ files for changes and regenerates them automatically
func Watch() error {
	if err := ensureTool(toolTempl); err != nil {
		return err
	}
	fmt.Println("👁️ Watching templ files for changes...")
	fmt.Println("Press Ctrl+C to stop watching")
	return runCmd(toolTempl, "generate", "--watch")
}

// CleanDev removes all development artifacts, generated files, and caches
func CleanDev() error {
	fmt.Println("🧹 Cleaning development artifacts...")
	_ = os.RemoveAll("tmp")

	fmt.Println("Cleaning node_modules and package files...")
	_ = os.RemoveAll("node_modules")
	_ = os.Remove("package.json")
	_ = os.Remove("package-lock.json")

	fmt.Println("Cleaning generated templ files...")
	if err := runCmd("find", ".", "-name", "*_templ.go", "-delete"); err != nil {
		fmt.Printf("Warning: find failed: %v\n", err)
		_ = runCmd("powershell", "-Command", "Get-ChildItem -Recurse -Filter '*_templ.go' | Remove-Item")
	}

	fmt.Println("Cleaning go build cache...")
	_ = runCmd("go", "clean", "-cache")
	return Clean()
}
