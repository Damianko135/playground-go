//go:build mage

package main

import (
	"runtime"
)

var Default = Build

// ─── Configuration ────────────────────────────────────────────────────────────

const (
	winBin   = "bin/main.exe"
	linuxBin = "bin/main"
	cmdPath  = "cmd/server/main.go"

	toolTempl        = "templ"
	toolAir          = "air"
	toolGolangciLint = "golangci-lint"
	toolGolds        = "golds"
	toolNpm          = "npm"
)

func outputBin() string {
	if runtime.GOOS == "windows" {
		return winBin
	}
	return linuxBin
}
