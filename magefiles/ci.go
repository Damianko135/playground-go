//go:build mage

package main

import (
	"fmt"

	"github.com/magefile/mage/mg"
)

// ─── Testing & CI ─────────────────────────────────────────────────────────────

// Test runs all tests
func Test() error {
	fmt.Println("🧪 Testing...")
	return runCmd("go", "test", "./...")
}

// Lint runs code quality checks
func Lint() error {
	if err := ensureTool(toolGolangciLint); err != nil {
		fmt.Println("⚠️ golangci-lint not found, using go vet")
		return lintFallback()
	}

	fmt.Println("🪄 Linting...")
	return runCmd(toolGolangciLint, "run")
}

// Fmt formats code
func Fmt() error {
	fmt.Println("🧹 Formatting...")
	return runCmd("gofmt", "-s", "-w", ".")
}

// Check runs quality checks
func Check() error {
	fmt.Println("🔍 Checking...")
	mg.SerialDeps(Fmt, Lint, Test)
	return nil
}

// CI runs complete pipeline
func CI() error {
	fmt.Println("🔁 CI...")
	mg.SerialDeps(Check, Build)
	return nil
}

// ─── Internal Functions ───────────────────────────────────────────────────────

func lintFallback() error {
	fmt.Println("🔍 Running go vet...")
	return runCmd("go", "vet", "./...")
}
