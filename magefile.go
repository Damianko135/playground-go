// //go:build mage

package main

// import (
// 	"fmt"
// 	"io"
// 	"os"
// 	"os/exec"
// 	"os/signal"
// 	"path/filepath"
// 	"runtime"
// 	"syscall"

// 	"github.com/magefile/mage/mg"
// 	"github.com/magefile/mage/sh"
// )

// var Default = Build

// // ─── Configuration ────────────────────────────────────────────────────────────

// const (
// 	winBin   = "bin/main.exe"
// 	linuxBin = "bin/main"
// 	cmdPath  = "cmd/server/main.go"

// 	toolTempl        = "templ"
// 	toolAir          = "air"
// 	toolGolangciLint = "golangci-lint"
// 	toolGolds        = "golds"
// 	toolNpm          = "npm"
// )

// func outputBin() string {
// 	if runtime.GOOS == "windows" {
// 		return winBin
// 	}
// 	return linuxBin
// }

// // ─── Utility Functions ────────────────────────────────────────────────────────

// func runCmd(name string, args ...string) error {
// 	fmt.Printf("Running: %s %v\n", name, args)
// 	return sh.Run(name, args...)
// }

// func ensureTool(tool string) error {
// 	_, err := exec.LookPath(tool)
// 	if err != nil {
// 		return fmt.Errorf("required tool '%s' not found in PATH", tool)
// 	}
// 	return nil
// }

// func copyFile(src, dst string) error {
// 	in, err := os.Open(src)
// 	if err != nil {
// 		return err
// 	}
// 	defer in.Close()

// 	out, err := os.Create(dst)
// 	if err != nil {
// 		return err
// 	}
// 	defer out.Close()

// 	if _, err := io.Copy(out, in); err != nil {
// 		return err
// 	}
// 	return out.Sync()
// }

// // ─── Build Targets ────────────────────────────────────────────────────────────

// func Clean() error {
// 	fmt.Println("🧹 Cleaning...")
// 	return os.RemoveAll("bin")
// }

// func Build() error {
// 	mg.Deps(Clean)
// 	fmt.Println("🔨 Building...")
// 	return runCmd("go", "build", "-o", outputBin(), cmdPath)
// }

// func Tags(tag string) error {
// 	fmt.Printf("🔨 Building with tag: %s\n", tag)
// 	return runCmd("go", "build", "-tags", tag, "-o", outputBin(), cmdPath)
// }

// func Install() error {
// 	dest := os.Getenv("GOBIN")
// 	if dest == "" {
// 		if runtime.GOOS == "windows" {
// 			return fmt.Errorf("GOBIN must be set on Windows")
// 		}
// 		dest = "/usr/local/bin"
// 	}
// 	fmt.Printf("📦 Installing to %s...\n", dest)
// 	binPath := outputBin()
// 	destPath := filepath.Join(dest, filepath.Base(binPath))
// 	return copyFile(binPath, destPath)
// }

// func Run() error {
// 	bin := outputBin()
// 	if _, err := os.Stat(bin); os.IsNotExist(err) {
// 		fmt.Println("⚠️ Binary missing — building first...")
// 		if err := Build(); err != nil {
// 			return err
// 		}
// 	}
// 	fmt.Println("🚀 Running...")
// 	return runCmd("./" + bin)
// }

// // ─── Development ──────────────────────────────────────────────────────────────
// func Dev() error {
// 	if err := ensureTool(toolTempl); err != nil {
// 		return err
// 	}
// 	if err := ensureTool(toolAir); err != nil {
// 		fmt.Println("⚠️ Air not found, falling back to basic dev mode")
// 		return DevBasic()
// 	}

// 	fmt.Println("🛠️ Generating templ files...")
// 	if err := runCmd(toolTempl, "generate"); err != nil {
// 		return err
// 	}

// 	tailwindPath := filepath.Join("node_modules", ".bin", "tailwindcss")
// 	if runtime.GOOS == "windows" {
// 		tailwindPath += ".cmd"
// 	}
// 	tailwindCmd := exec.Command(tailwindPath, "-i", "./internal/app.css", "-o", "./static/index.css", "--minify")
// 	tailwindCmd.Stdout = os.Stdout
// 	tailwindCmd.Stderr = os.Stderr
// 	tailwindCmd.Env = os.Environ()
// 	tailwindCmd.Dir = "." // Run from project root

// 	fmt.Println("📦 Running Tailwind build with npx...")
// 	if err := tailwindCmd.Start(); err != nil {
// 		return fmt.Errorf("failed to start tailwind watcher via npx: %w", err)
// 	}
// 	fmt.Printf("✅ Tailwind watcher started (pid: %d)\n", tailwindCmd.Process.Pid)

// 	airCmd := exec.Command(toolAir, "-c", "../../.air.toml")
// 	airCmd.Dir = "cmd/server"
// 	airCmd.Env = append(os.Environ(), "GO_ENV=development", "PORT=3000")
// 	airCmd.Stdout = os.Stdout
// 	airCmd.Stderr = os.Stderr
// 	airCmd.Stdin = os.Stdin

// 	fmt.Println("🚀 Starting development server with hot reload...")
// 	fmt.Println("📍 http://localhost:3000")

// 	sigs := make(chan os.Signal, 1)
// 	done := make(chan error, 1)
// 	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

// 	go func() {
// 		done <- airCmd.Run()
// 	}()

// 	select {
// 	case sig := <-sigs:
// 		fmt.Printf("\nReceived signal: %s, shutting down...\n", sig)
// 		_ = tailwindCmd.Process.Kill()
// 		_ = airCmd.Process.Kill()
// 		return nil
// 	case err := <-done:
// 		_ = tailwindCmd.Process.Kill()
// 		return err
// 	}
// }

// func DevBasic() error {
// 	fmt.Println("🛠️ Generating templ files...")
// 	if err := runCmd(toolTempl, "generate"); err != nil {
// 		return err
// 	}

// 	fmt.Println("🔄 Running in basic dev mode...")
// 	cmd := exec.Command("go", "run", ".")
// 	cmd.Dir = "cmd/server"
// 	cmd.Env = append(os.Environ(), "GO_ENV=development", "PORT=3000")
// 	cmd.Stdout = os.Stdout
// 	cmd.Stderr = os.Stderr
// 	cmd.Stdin = os.Stdin
// 	return cmd.Run()
// }

// func Watch() error {
// 	if err := ensureTool(toolTempl); err != nil {
// 		return err
// 	}
// 	fmt.Println("👁️ Watching templ files for changes...")
// 	fmt.Println("Press Ctrl+C to stop watching")
// 	return runCmd(toolTempl, "generate", "--watch")
// }

// func DevClean() error {
// 	fmt.Println("🧹 Cleaning development artifacts...")
// 	_ = os.RemoveAll("tmp")

// 	fmt.Println("Cleaning generated templ files...")
// 	if err := runCmd("find", ".", "-name", "*_templ.go", "-delete"); err != nil {
// 		fmt.Printf("Warning: find failed: %v\n", err)
// 		_ = runCmd("powershell", "-Command", "Get-ChildItem -Recurse -Filter '*_templ.go' | Remove-Item")
// 	}

// 	fmt.Println("Cleaning go build cache...")
// 	_ = runCmd("go", "clean", "-cache")
// 	return Clean()
// }

// // ─── Testing & CI ─────────────────────────────────────────────────────────────

// func Test() error {
// 	fmt.Println("🧪 Running tests...")
// 	return runCmd("go", "test", "./...")
// }

// func Lint() error {
// 	if err := ensureTool(toolGolangciLint); err != nil {
// 		fmt.Println("Install: https://golangci-lint.run/usage/install/")
// 		return err
// 	}
// 	fmt.Println("🪄 Running linter...")
// 	return runCmd(toolGolangciLint, "run")
// }

// func Format() error {
// 	fmt.Println("🧹 Formatting code...")
// 	return runCmd("gofmt", "-s", "-w", ".")
// }

// func CI() error {
// 	fmt.Println("🔁 Running CI steps...")
// 	mg.SerialDeps(Lint, Test, Build)
// 	return nil
// }

// // ─── Docs Namespace ──────────────────────────────────────────────────────────

// type Docs mg.Namespace

// func (Docs) Gen() error {
// 	if err := ensureTool("gomarkdoc"); err != nil {
// 		return err
// 	}
// 	fmt.Println("📚 Generating Markdown docs with gomarkdoc...")
// 	return runCmd("gomarkdoc", "-u", "-o", "DOCS.md", "./...")
// }

// // ─── Tool Management ─────────────────────────────────────────────────────────

// func DevTools() error {
// 	tools := []struct {
// 		name string
// 		pkg  string
// 	}{
// 		{toolTempl, "github.com/a-h/templ/cmd/templ@latest"},
// 		{toolAir, "github.com/cosmtrek/air@latest"},
// 		{toolGolangciLint, "github.com/golangci/golangci-lint/cmd/golangci-lint@latest"},
// 	}
// 	fmt.Println("🛠️ Installing development tools...")
// 	for _, tool := range tools {
// 		fmt.Printf("Installing %s...\n", tool.name)
// 		if err := runCmd("go", "install", tool.pkg); err != nil {
// 			fmt.Printf("❌ Failed to install %s: %v\n", tool.name, err)
// 			continue
// 		}
// 		fmt.Printf("✅ %s installed successfully\n", tool.name)
// 	}
// 	return nil
// }

// // ─── Misc ─────────────────────────────────────────────────────────────────────

// func Srv() error {
// 	fmt.Println("🌀 Running server package...")
// 	return runCmd("go", "run", "./cmd/server")
// }

// func Hot() error {
// 	if err := ensureTool(toolAir); err != nil {
// 		return err
// 	}
// 	fmt.Println("♻️ Running with Air (hot reload)...")
// 	return runCmd(toolAir)
// }

// func Help() {
// 	fmt.Println(`Available targets:
//   clean       - Remove build artifacts
//   build       - Build the main binary
//   tags        - Build with custom tags (usage: mage tags -tag=debug)
//   run         - Run the built binary
//   dev         - Run development mode with hot reload (or fallback)
//   devbasic    - Run development mode without hot reload
//   hot         - Run hot reload using air
//   srv         - Run server package
//   test        - Run tests
//   lint        - Run golangci-lint
//   install     - Install binary to GOBIN or /usr/local/bin
//   docs:gen    - Generate documentation
//   ci          - Run CI steps: lint, test, build
//   format      - Format code with gofmt
//   watch       - Watch templ files and regenerate on changes
//   devtools    - Install development tools
//   devclean    - Clean development artifacts and caches
//   help        - Show this help message
// `)
// }
