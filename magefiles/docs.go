//go:build mage

package main

import (
	"fmt"

	"github.com/magefile/mage/mg"
)

// ─── Docs Namespace ──────────────────────────────────────────────────────────

type Docs mg.Namespace

// Gen generates Markdown documentation from Go code using gomarkdoc
func (Docs) Gen() error {
	if err := ensureTool("gomarkdoc"); err != nil {
		return err
	}
	fmt.Println("📚 Generating Markdown docs with gomarkdoc...")
	return runCmd("gomarkdoc", "-u", "-o", "DOCS.md", "./...")
}
