//go:build mage

package main

import (
	"fmt"
)

// ─── Tool Management ─────────────────────────────────────────────────────────

// Tools installs development tools
func Tools() error {
	tools := []struct {
		name string
		pkg  string
	}{
		{toolTempl, "github.com/a-h/templ/cmd/templ@latest"},
		{toolAir, "github.com/air-verse/air@latest"},
		{toolGolangciLint, "github.com/golangci/golangci-lint/cmd/golangci-lint@latest"},
		{toolMage, "github.com/magefile/mage@latest"},
		{toolGomarkdoc, "github.com/princjef/gomarkdoc/cmd/gomarkdoc@latest"},
	}

	fmt.Println("🛠️ Installing tools...")

	var failed []string
	for _, tool := range tools {
		fmt.Printf("📦 %s...\n", tool.name)
		if err := runCmd("go", "install", tool.pkg); err != nil {
			fmt.Printf("❌ %s failed\n", tool.name)
			failed = append(failed, tool.name)
			continue
		}
		fmt.Printf("✅ %s\n", tool.name)
	}

	if len(failed) > 0 {
		return fmt.Errorf("failed: %v", failed)
	}

	fmt.Println("🎉 Done!")
	return nil
}
