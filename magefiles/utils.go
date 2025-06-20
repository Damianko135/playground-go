//go:build mage

package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/magefile/mage/sh"
)

// ─── Utility Functions ────────────────────────────────────────────────────────

func runCmd(name string, args ...string) error {
	fmt.Printf("Running: %s %v\n", name, args)
	return sh.Run(name, args...)
}

func ensureTool(tool string) error {
	_, err := exec.LookPath(tool)
	if err != nil {
		return fmt.Errorf("required tool '%s' not found in PATH", tool)
	}
	return nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return out.Sync()
}

