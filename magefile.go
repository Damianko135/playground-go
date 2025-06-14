//go:build mage

package main

import (
	"fmt"
	"os"
	"os/exec"
)

// Default task: build the project.
var Default = Build

// Build compiles the Go application into ./bin/main.exe (Windows-compatible).
func Build() error {
	fmt.Println("🔨 Building the application...")
	cmd := exec.Command("go", "build", "-o", "./bin/main.exe", "cmd/server/main.go")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Run starts the compiled application.
func Run() error {
	fmt.Println("🚀 Running the application...")
	cmd := exec.Command("./bin/main.exe")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

// Dev runs templ generate and starts the app using `go run` (no build).
func Dev() error {
	fmt.Println("🛠️  Running in development mode...")

	// Generate Templ components
	if err := exec.Command("templ", "generate").Run(); err != nil {
		fmt.Println("Templ generation failed:", err)
		return err
	}

	cmd := exec.Command("go", "run", "cmd/server/main.go")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

// Clean deletes the ./bin directory.
func Clean() error {
	fmt.Println("🧹 Cleaning build artifacts...")
	return os.RemoveAll("./bin")
}

// Test runs all Go tests in the project.
func Test() error {
	fmt.Println("🧪 Running tests...")
	cmd := exec.Command("go", "test", "./...")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
