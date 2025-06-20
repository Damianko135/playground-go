//go:build mage

package main

import (
	"fmt"
)

// ─── Misc ─────────────────────────────────────────────────────────────────────

// Hot runs the server with Air hot reload, falls back to basic run if Air is missing
func Hot() error {
	if err := ensureTool(toolAir); err != nil {
		fmt.Println("⚠️ Air not found, falling back to basic server run")
		return srv()
	}
	fmt.Println("♻️ Running with Air (hot reload)...")
	return runCmd(toolAir)
}

// srv is an internal fallback function for basic server running
func srv() error {
	fmt.Println("🌀 Running server package...")
	return runCmd("go", "run", "./cmd/server")
}

// Help displays information about available mage targets
func Help() {
	fmt.Println(`Available targets:
  clean       - Remove build artifacts
  build       - Build the main binary
  tags        - Build with custom tags (usage: mage tags -tag=debug)
  run         - Run the built binary
  dev         - Run development mode (with hot reload, falls back to basic mode)
  hot         - Run server with hot reload (falls back to basic run)
  test        - Run tests
  lint        - Run linter (golangci-lint, falls back to go vet)
  install     - Install binary to GOBIN or /usr/local/bin
  docs:gen    - Generate documentation
  ci          - Run CI steps: lint, test, build
  fmt         - Format code with gofmt
  watch       - Watch templ files and regenerate on changes
  tools       - Install development tools
  cleandev    - Clean development artifacts and caches
  help        - Show this help message
`)
}
