//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"os"
	"os/exec"
)

const magePkg = "github.com/magefile/mage@latest"

func main() {
	fmt.Println("🔍 Checking for mage...")

	if _, err := exec.LookPath("mage"); err != nil {
		fmt.Println("📦 'mage' not found. Installing...")
		installMage := exec.Command("go", "install", magePkg)
		installMage.Stdout = os.Stdout
		installMage.Stderr = os.Stderr
		installMage.Env = os.Environ()

		if err := installMage.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "❌ Failed to install mage: %v\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Println("✅ 'mage' is already installed.")
	}

	fmt.Println("🚀 Running 'mage tools'...")

	cmd := exec.Command("mage", "tools")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Failed to run 'mage tools': %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ All dev tools installed.")
}
