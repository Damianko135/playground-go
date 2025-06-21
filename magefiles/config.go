//go:build mage

package main

import (
	"runtime"
)

var Default = Help

// ─── Configuration ────────────────────────────────────────────────────────────

const (
	// Binary paths
	winBin   = "bin/main.exe"
	linuxBin = "bin/main"
	cmdPath  = "cmd/server/main.go"

	// Development tools
	toolTempl        = "templ"
	toolAir          = "air"
	toolGolangciLint = "golangci-lint"
	toolGolds        = "golds"
	toolGomarkdoc    = "gomarkdoc"
	toolNpm          = "npm"
	toolMage         = "mage"

	// Development settings
	devPort = "3000"
	devEnv  = "GO_ENV=development"
)

func outputBin() string {
	if runtime.GOOS == "windows" {
		return winBin
	}
	return linuxBin
}
