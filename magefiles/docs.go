//go:build mage

package main

import (
	"fmt"

	"github.com/magefile/mage/mg"
)

// ─── Documentation ────────────────────────────────────────────────────────────

type Docs mg.Namespace

// Gen generates documentation
func (Docs) Gen() error {
	if err := ensureTool(toolGomarkdoc); err != nil {
		return fmt.Errorf("install gomarkdoc: go install github.com/princjef/gomarkdoc/cmd/gomarkdoc@latest")
	}

	fmt.Println("📚 Generating docs...")
	return runCmd(toolGomarkdoc, "-u", "-o", "DOCS.md", "./...")
}

// Serve starts documentation server
func (Docs) Serve() error {
	if err := ensureTool(toolGolds); err != nil {
		return fmt.Errorf("install golds: go install go101.org/golds@latest")
	}

	fmt.Println("📖 Starting docs server...")
	fmt.Println("📍 http://localhost:8080")
	return runCmd(toolGolds, "-port=8080", "./cmd/...", "./internal/...", "./views")
}
