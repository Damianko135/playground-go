//go:build mage

package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/magefile/mage/sh"
)

// ─── Utility Functions ────────────────────────────────────────────────────────

// runCmd executes a command with arguments and prints the command being run
func runCmd(name string, args ...string) error {
	fmt.Printf("▶️ Running: %s %v\n", name, args)
	return sh.Run(name, args...)
}

// runCmdQuiet executes a command without printing it (for internal operations)
func runCmdQuiet(name string, args ...string) error {
	return sh.Run(name, args...)
}

// ensureTool checks if a required tool is available in PATH
func ensureTool(tool string) error {
	_, err := exec.LookPath(tool)
	if err != nil {
		return fmt.Errorf("required tool '%s' not found in PATH", tool)
	}
	return nil
}

// copyFile copies a file from src to dst with proper error handling
func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file %s: %w", src, err)
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file %s: %w", dst, err)
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}

	if err := dstFile.Sync(); err != nil {
		return fmt.Errorf("failed to sync file: %w", err)
	}

	// Copy file permissions
	srcInfo, err := srcFile.Stat()
	if err != nil {
		return fmt.Errorf("failed to get source file info: %w", err)
	}

	if err := os.Chmod(dst, srcInfo.Mode()); err != nil {
		return fmt.Errorf("failed to set file permissions: %w", err)
	}

	return nil
}

// fileExists checks if a file exists
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// removeIfExists removes a file or directory if it exists, ignoring "not found" errors
func removeIfExists(path string) error {
	err := os.RemoveAll(path)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove %s: %w", path, err)
	}
	return nil
}

// ─── Environment Utilities ────────────────────────────────────────────────────

// getEnvWithDefault returns environment variable value or default
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// setDefaultEnv sets environment variable only if not already set
func setDefaultEnv(key, value string) {
	if os.Getenv(key) == "" {
		os.Setenv(key, value)
	}
}

// loadEnvFromFile loads environment variables from a file
func loadEnvFromFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse KEY=VALUE format
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Remove quotes if present
		if len(value) >= 2 {
			if (strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"")) ||
				(strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'")) {
				value = value[1 : len(value)-1]
			}
		}

		// Only set if not already set (environment variables take precedence)
		if os.Getenv(key) == "" {
			os.Setenv(key, value)
		}
	}

	return nil
}
