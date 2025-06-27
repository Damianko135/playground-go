//go:build mage

package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// ─── Docker Targets ──────────────────────────────────────────────────────────

// DockerBuild builds the Docker image for the application
func DockerBuild() error {
	fmt.Println("🐳 Building Docker image...")
	return sh.RunWith(map[string]string{"PWD": getProjectRoot()}, "docker", "build", "-t", "playground-go:latest", ".")
}

// DockerBuildDev builds the Docker image with development configuration
func DockerBuildDev() error {
	fmt.Println("🐳 Building Docker image for development...")
	return sh.RunWith(map[string]string{"PWD": getProjectRoot()}, "docker", "build", "-t", "playground-go:dev", "--build-arg", "CONFIGURATION=simple", ".")
}

// DockerBuildProd builds the Docker image with production optimizations
func DockerBuildProd() error {
	fmt.Println("🐳 Building optimized Docker image for production...")
	return sh.RunWith(map[string]string{"PWD": getProjectRoot()}, "docker", "build", "-t", "playground-go:prod", "--build-arg", "CONFIGURATION=optimized", ".")
}

// DockerRun runs the application in a Docker container (development mode)
func DockerRun() error {
	mg.Deps(DockerBuild)
	fmt.Println("🚀 Running application in Docker container...")

	args := []string{
		"run",
		"--rm",
		"-it",
		"-p", "8080:8080",
		"-v", fmt.Sprintf("%s:/app", getProjectRoot()),
		"-w", "/app",
		"-e", "ENV=development",
		"playground-go:latest",
	}

	return sh.Run("docker", args...)
}

// DockerRunProd runs the application in production mode
func DockerRunProd() error {
	mg.Deps(DockerBuildProd)
	fmt.Println("🚀 Running application in Docker container (production)...")

	args := []string{
		"run",
		"--rm",
		"-it",
		"-p", "3000:8080",
		"-e", "ENV=production",
		"playground-go:prod",
	}

	return sh.Run("docker", args...)
}

// DockerDev starts the development environment using docker-compose
func DockerDev() error {
	fmt.Println("🐳 Starting development environment with docker-compose...")
	return sh.RunWith(map[string]string{"PWD": getProjectRoot()}, "docker", "compose", "-f", "compose.yaml", "up", "app_dev", "--build")
}

// DockerProd starts the production environment using docker-compose
func DockerProd() error {
	fmt.Println("🐳 Starting production environment with docker-compose...")
	return sh.Run("docker", "compose", "-f", "compose.yaml", "up", "app_prod", "--build", "-d")
}

// DockerStop stops all running Docker containers for this project
func DockerStop() error {
	fmt.Println("🛑 Stopping Docker containers...")
	return sh.Run("docker", "compose", "-f", "compose.yaml", "down")
}

// DockerLogs shows logs from the Docker containers
func DockerLogs() error {
	fmt.Println("📋 Showing Docker container logs...")
	return sh.Run("docker", "compose", "-f", "compose.yaml", "logs", "-f")
}

// DockerClean removes Docker images and containers for this project
func DockerClean() error {
	fmt.Println("🧹 Cleaning up Docker resources...")

	// Stop containers first
	sh.Run("docker", "compose", "-f", "compose.yaml", "down")

	// Remove images
	sh.Run("docker", "rmi", "playground-go:latest", "playground-go:dev", "playground-go:prod")

	// Clean up dangling images
	return sh.Run("docker", "image", "prune", "-f")
}

// DockerShell opens an interactive shell in the development container
func DockerShell() error {
	mg.Deps(DockerBuild)
	fmt.Println("🐚 Opening shell in Docker container...")

	args := []string{
		"run",
		"--rm",
		"-it",
		"-v", fmt.Sprintf("%s:/app", getProjectRoot()),
		"-w", "/app",
		"playground-go:latest",
		"sh",
	}

	return sh.Run("docker", args...)
}

// DockerStatus checks if Docker is running and shows container status
func DockerStatus() error {
	fmt.Println("🔍 Checking Docker status...")

	// Check if Docker daemon is running
	if err := sh.Run("docker", "info"); err != nil {
		fmt.Println("❌ Docker daemon is not running")
		return err
	}

	fmt.Println("✅ Docker daemon is running")

	// Show running containers for this project
	fmt.Println("📋 Docker containers for this project:")
	return sh.Run("docker", "compose", "-f", "compose.yaml", "ps")
}

// DockerRestart restarts the development environment
func DockerRestart() error {
	fmt.Println("🔄 Restarting development environment...")

	// Stop containers
	sh.Run("docker", "compose", "-f", "compose.yaml", "down")

	// Start them again
	return sh.Run("docker", "compose", "-f", "compose.yaml", "up", "app_dev", "--build", "-d")
}

// getProjectRoot returns the absolute path to the project root
func getProjectRoot() string {
	// Start from the current working directory (which should be magefiles)
	wd, err := os.Getwd()
	if err != nil {
		return "."
	}

	// If we're in the magefiles directory, go up one level to project root
	if filepath.Base(wd) == "magefiles" {
		return filepath.Dir(wd)
	}

	// Otherwise assume we're already in project root
	return wd
}
