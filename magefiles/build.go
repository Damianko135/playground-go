//go:build mage

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/magefile/mage/mg"
)

// ─── Build Targets ────────────────────────────────────────────────────────────

// Clean removes build artifacts and temporary files
func Clean() error {
	fmt.Println("🧹 Cleaning...")

	dirs := []string{"bin", "cmd/server/tmp", "node_modules"}
	for _, dir := range dirs {
		if err := os.RemoveAll(dir); err != nil && !os.IsNotExist(err) {
			fmt.Printf("Warning: failed to remove %s: %v\n", dir, err)
		}
	}
	return nil
}

// Build compiles the application
func Build() error {
	mg.Deps(Clean)
	fmt.Println("🔨 Building...")
	return runCmd("go", "build", "-o", outputBin(), cmdPath)
}

// Run executes the built binary
func Run() error {
	bin := outputBin()
	if _, err := os.Stat(bin); os.IsNotExist(err) {
		fmt.Println("⚠️ Binary not found, building...")
		if err := Build(); err != nil {
			return fmt.Errorf("build failed: %w", err)
		}
	}

	fmt.Println("🚀 Running...")
	return runCmd("./" + bin)
}

// Install copies binary to GOBIN or system path
func Install() error {
	mg.Deps(Build)

	dest := os.Getenv("GOBIN")
	if dest == "" {
		if runtime.GOOS == "windows" {
			return fmt.Errorf("GOBIN must be set on Windows")
		}
		dest = "/usr/local/bin"
	}

	fmt.Printf("📦 Installing to %s...\n", dest)
	binPath := outputBin()
	destPath := filepath.Join(dest, filepath.Base(binPath))
	return copyFile(binPath, destPath)
}

// Gen generates templ files
func Gen() error {
	fmt.Println("🔧 Generating...")
	if err := ensureTool(toolTempl); err != nil {
		return err
	}
	return runCmd(toolTempl, "generate")
}

// Prod builds and runs the application in production mode
func Prod() error {
	fmt.Println("🚀 Production mode...")
	
	// Load environment from .env file if it exists
	loadEnvFile()
	
	// Generate templates
	mg.Deps(Gen)
	
	// Build CSS for production
	fmt.Println("📦 Building production CSS...")
	if err := buildProductionCSS(); err != nil {
		return fmt.Errorf("CSS build failed: %w", err)
	}
	
	// Build application with optimizations
	fmt.Println("🔨 Building for production...")
	if err := buildProduction(); err != nil {
		return fmt.Errorf("production build failed: %w", err)
	}
	
	// Show configuration
	showProductionConfig()
	
	// Run in production mode
	fmt.Println("🌐 Starting production server...")
	return runProductionServer()
}

// ─── Production Helper Functions ──────────────────────────────────────────────

// loadEnvFile loads environment variables from .env file if it exists
func loadEnvFile() {
	envFile := ".env"
	if fileExists(envFile) {
		fmt.Printf("📄 Loading environment from %s...\n", envFile)
		if err := loadEnvFromFile(envFile); err != nil {
			fmt.Printf("⚠️ Warning: failed to load %s: %v\n", envFile, err)
		}
	} else {
		fmt.Println("📄 No .env file found, using system environment variables")
		fmt.Println("💡 Tip: Copy .env.example to .env to customize configuration")
	}
}

// showProductionConfig displays the current configuration
func showProductionConfig() {
	port := getEnvWithDefault("PORT", "8080")
	host := getEnvWithDefault("HOST", "0.0.0.0")
	env := getEnvWithDefault("GO_ENV", "production")
	
	fmt.Println("⚙️ Production Configuration:")
	fmt.Printf("   🌐 Server: http://%s:%s\n", host, port)
	fmt.Printf("   🏷️ Environment: %s\n", env)
	fmt.Printf("   🔧 Debug: %s\n", getEnvWithDefault("DEBUG", "false"))
}

// buildProductionCSS builds optimized CSS for production
func buildProductionCSS() error {
	// Install Tailwind if not present
	if err := runCmd("npm", "install", "-D", "tailwindcss@latest", "@tailwindcss/cli@latest"); err != nil {
		return err
	}
	
	// Build minified CSS
	return runCmd("npx", "@tailwindcss/cli",
		"-i", "./internal/app.css",
		"-o", "./static/index.css",
		"--minify")
}

// buildProduction builds the application with production optimizations
func buildProduction() error {
	mg.Deps(Clean)
	
	// Build with optimizations: disable debug info, strip symbols, optimize for size
	return runCmd("go", "build",
		"-ldflags", "-s -w", // Strip debug info and symbol table
		"-trimpath",         // Remove file system paths from executable
		"-o", outputBin(),
		cmdPath)
}

// runProductionServer runs the server in production mode
func runProductionServer() error {
	bin := outputBin()
	
	// Set default production environment variables if not already set
	setDefaultEnv("GO_ENV", "production")
	setDefaultEnv("PORT", "8080")
	setDefaultEnv("HOST", "0.0.0.0")
	setDefaultEnv("DEBUG", "false")
	setDefaultEnv("ENABLE_PROFILING", "false")
	
	fmt.Printf("📍 Server will be available at: http://%s:%s\n", 
		getEnvWithDefault("HOST", "0.0.0.0"), 
		getEnvWithDefault("PORT", "8080"))
	
	// Run with current environment
	return runCmd("./"+bin)
}
