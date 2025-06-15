//go:build mage

package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/magefile/mage/mg"
)

// Default task: build the project.
var Default = Build

// Binary names
const (
	winBin  = "bin/main.exe"
	linuxBin = "bin/main"
)

// Build compiles the Go application into ./bin/
func Build() error {
	fmt.Println("🧹 Cleaning up previous builds...")
	_ = os.RemoveAll("./bin")
	_ = os.MkdirAll("./bin", 0755)

	fmt.Println("🔨 Building the application...")
	output := linuxBin
	if runtime.GOOS == "windows" {
		output = winBin
	}

	cmd := exec.Command("go", "build", "-o", output, "cmd/server/main.go")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// BuildWithTags builds with optional tags (e.g., mage buildWithTags debug).
func BuildWithTags(tags string) error {
	fmt.Printf("🔨 Building with tags: %s\n", tags)
	output := linuxBin
	if runtime.GOOS == "windows" {
		output = winBin
	}

	cmd := exec.Command("go", "build", "-tags", tags, "-o", output, "cmd/server/main.go")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Run starts the compiled binary.
func Run() error {
	fmt.Println("🚀 Running the application...")

	bin := linuxBin
	if runtime.GOOS == "windows" {
		bin = winBin
	}

	if _, err := os.Stat(bin); os.IsNotExist(err) {
		fmt.Println("Binary not found, building first...")
		if err := Build(); err != nil {
			return fmt.Errorf("build failed: %w", err)
		}
	}

	cmd := exec.Command("./" + bin)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

// Dev mode: generate templ files and run uncompiled app
func Dev() error {
	fmt.Println("🛠️  Running in development mode...")

	if err := exec.Command("templ", "generate").Run(); err != nil {
		fmt.Println("Templ generation failed:", err)
		return err
	}

	cmd := exec.Command("go", "run", "cmd/server/main.go")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

// Install copies the binary to $GOBIN or /usr/local/bin
func Install() error {
	fmt.Println("📦 Installing binary...")

	dest := os.Getenv("GOBIN")
	if dest == "" {
		dest = "/usr/local/bin"
	}
	binary := linuxBin
	if runtime.GOOS == "windows" {
		binary = winBin
		dest = os.Getenv("GOBIN") // user-defined GOBIN is a must
	}
	if dest == "" {
		return fmt.Errorf("GOBIN not set, and default location invalid on Windows")
	}

	fmt.Printf("📂 Installing to %s\n", dest)
	cmd := exec.Command("cp", binary, dest)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Clean removes bin and generated files
func Clean() error {
	fmt.Println("🧹 Cleaning build artifacts...")
	return os.RemoveAll("./bin")
}

// Tidy runs go mod tidy
func Tidy() error {
	fmt.Println("🧹 Tidying dependencies...")
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Test runs all unit tests
func Test() error {
	fmt.Println("🧪 Running tests...")
	cmd := exec.Command("go", "test", "./...")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Lint runs golangci-lint
func Lint() error {
	fmt.Println("🪄 Checking code style...")
	if err := exec.Command("golangci-lint", "--version").Run(); err != nil {
		fmt.Println("golangci-lint is not installed. Install: https://golangci-lint.run/usage/install/")
		return err
	}
	cmd := exec.Command("golangci-lint", "run")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Generate runs all code generation tools
func Generate() error {
	fmt.Println("⚙️  Running go generate...")
	cmd := exec.Command("go", "generate", "./...")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	fmt.Println("🧬 Running templ generate...")
	cmd = exec.Command("templ", "generate")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Docker builds the Docker image
func Docker() error {
	fmt.Println("🐳 Building Docker image...")
	cmd := exec.Command("docker", "build", "-t", "myapp:latest", ".")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// CI runs lint, test, and build — for CI pipelines
func CI() error {
	fmt.Println("🔁 Running CI tasks...")
	if err := Lint(); err != nil {
		return err
	}
	if err := Test(); err != nil {
		return err
	}
	return Build()
}

// Docs task group
type Docs mg.Namespace

// Generate markdown docs
func (Docs) Generate() error {
	fmt.Println("📚 Generating static documentation with golds...")

	cmd := exec.Command("golds "," -out=docs.md " ," -gen "," -dir=generated "," -nouses "," std")

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}



