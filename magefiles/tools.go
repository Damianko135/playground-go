//go:build mage

package main

import (
	"fmt"
)

// ─── Tool Management ─────────────────────────────────────────────────────────

// Tools installs all required development tools (templ, air, golangci-lint)
func Tools() error {
	tools := []struct {
		name string
		pkg  string
	}{
		{toolTempl, "github.com/a-h/templ/cmd/templ@latest"},
		{toolAir, "github.com/cosmtrek/air@latest"},
		{toolGolangciLint, "github.com/golangci/golangci-lint/cmd/golangci-lint@latest"},
	}
	fmt.Println("🛠️ Installing development tools...")
	for _, tool := range tools {
		fmt.Printf("Installing %s...\n", tool.name)
		if err := runCmd("go", "install", tool.pkg); err != nil {
			fmt.Printf("❌ Failed to install %s: %v\n", tool.name, err)
			continue
		}
		fmt.Printf("✅ %s installed successfully\n", tool.name)
	}
	return nil
}
