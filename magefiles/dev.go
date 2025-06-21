//go:build mage

package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

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
	cmd := exec.Command("go", "run", ".")
	cmd.Dir = "cmd/server"
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
	cmd := exec.Command(toolAir, "-c", "../../.air.toml")
	cmd.Dir = "cmd/server"
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
		cleanup()
		return nil
	case err := <-done:
		cleanup()
		return err
	}
}

func cleanupTailwind(cmd *exec.Cmd) {
	if cmd != nil && cmd.Process != nil {
		_ = cmd.Process.Kill()
	}
	fmt.Println("🧹 Cleaning up...")
	_ = os.RemoveAll("node_modules")
	_ = os.Remove("package.json")
	_ = os.Remove("package-lock.json")
}
