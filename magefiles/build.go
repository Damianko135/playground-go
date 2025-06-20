//go:build mage

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/magefile/mage/mg"
)

// ─── Build Targets ────────────────────────────────────────────────────────────

// Clean removes all build artifacts and the bin directory
func Clean() error {
	fmt.Println("🧹 Cleaning...")
	return os.RemoveAll("bin")
}

// Build compiles the application binary after cleaning
func Build() error {
	mg.Deps(Clean)
	fmt.Println("🔨 Building...")
	return runCmd("go", "build", "-o", outputBin(), cmdPath)
}

// Tags builds the application with custom build tags
func Tags(tag string) error {
	fmt.Printf("🔨 Building with tag: %s\n", tag)
	return runCmd("go", "build", "-tags", tag, "-o", outputBin(), cmdPath)
}

// Install copies the built binary to GOBIN or /usr/local/bin
func Install() error {
	dest := os.Getenv("GOBIN")
	if dest == "" {
		if runtime.GOOS == "windows" {
			return fmt.Errorf("GOBIN must be set on Windows")
		}
		dest = "/usr/local/bin"
	}
	fmt.Printf("📦 Installing to %s...\n", dest)
	binPath := outputBin()
	destPath := filepath.Join(dest, filepath.Base(binPath))
	return copyFile(binPath, destPath)
}

// Run executes the built binary, building it first if necessary
func Run() error {
	bin := outputBin()
	if _, err := os.Stat(bin); os.IsNotExist(err) {
		fmt.Println("⚠️ Binary missing — building first...")
		if err := Build(); err != nil {
			return err
		}
	}
	fmt.Println("🚀 Running...")
	return runCmd("./" + bin)
}
