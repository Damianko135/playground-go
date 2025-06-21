//go:build mage

package main

import (
	"fmt"
	"os"

	"github.com/magefile/mage/mg"
)

// ─── Miscellaneous ────────────────────────────────────────────────────────────

// Help shows available targets
func Help() {
	fmt.Println(`🚀 Mage Targets:

📦 Build:
  clean    - Remove build artifacts
  build    - Build application
  run      - Run application
  prod     - Build and run in production mode
  install  - Install to system

🛠️ Development:
  dev      - Start dev environment
  gen      - Generate templates
  watch    - Watch template files

🧪 Quality:
  test     - Run tests
  lint     - Run linter
  fmt      - Format code
  check    - Run all checks
  ci       - Full CI pipeline

🔧 Tools:
  tools    - Install dev tools

📚 Docs:
  docs:gen   - Generate docs
  docs:serve - Serve docs

Use 'mage -l' to list all targets.`)
}

// CleanAll removes all artifacts and caches
func CleanAll() error {
	fmt.Println("🧹 Deep clean...")

	mg.Deps(Clean)

	// Clean dev artifacts
	dirs := []string{"tmp", "node_modules"}
	files := []string{"package.json", "package-lock.json"}

	for _, dir := range dirs {
		_ = os.RemoveAll(dir)
	}
	for _, file := range files {
		_ = os.Remove(file)
	}

	// Clean generated files
	_ = runCmd("find", ".", "-name", "*_templ.go", "-delete")

	// Clean go cache
	_ = runCmd("go", "clean", "-cache")

	return nil
}
