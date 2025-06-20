//go:build mage

package main

import (
	"fmt"

	"github.com/magefile/mage/mg"
)

// ─── Testing & CI ─────────────────────────────────────────────────────────────

// Test runs all Go tests in the project
func Test() error {
	fmt.Println("🧪 Running tests...")
	return runCmd("go", "test", "./...")
}

// Lint runs golangci-lint for code quality checks, falls back to go vet if not available
func Lint() error {
	if err := ensureTool(toolGolangciLint); err != nil {
		fmt.Println("⚠️ golangci-lint not found, falling back to go vet")
		fmt.Println("For better linting, install: https://golangci-lint.run/usage/install/")
		return lintBasic()
	}
	fmt.Println("🪄 Running linter...")
	return runCmd(toolGolangciLint, "run")
}

// lintBasic is an internal fallback function using go vet
func lintBasic() error {
	fmt.Println("🔍 Running go vet...")
	return runCmd("go", "vet", "./...")
}

// Fmt formats all Go code using gofmt
func Fmt() error {
	fmt.Println("🧹 Formatting code...")
	return runCmd("gofmt", "-s", "-w", ".")
}

// CI runs the complete CI pipeline: lint, test, and build
func CI() error {
	fmt.Println("🔁 Running CI steps...")
	mg.SerialDeps(Lint, Test, Build)
	return nil
}
