# playground-go

A Go web application built with Echo framework, Templ for templating, and Tailwind CSS for styling.

## Features

- 🚀 **Echo Framework** - High performance, minimalist Go web framework
- 🎨 **Templ** - Type-safe HTML templating for Go
- 💨 **Tailwind CSS** - Utility-first CSS framework
- 🔥 **Hot Reload** - Development server with automatic reloading
- 🛠️ **Mage** - Build automation tool for Go

## Prerequisites

- Go 1.24.4 or later
- Node.js and npx (for Tailwind CSS)
- [Mage](https://magefile.org/) - Install with: `go install github.com/magefile/mage@latest`
- [Templ](https://templ.guide/) - Install with: `go install github.com/a-h/templ/cmd/templ@latest`

## Project Structure

```
.
├── cmd/                 # Application entry points
│   ├── cli/            # CLI application
│   └── server/         # Web server
├── internal/           # Private application code
│   ├── utils/          # Utility functions
│   └── app.css         # Main Tailwind CSS input file
├── magefiles/          # Mage build scripts
├── static/             # Static assets (compiled CSS, images)
│   └── index.css       # Compiled Tailwind CSS output
├── views/              # Templ templates (.templ files)
├── .air.toml           # Air configuration for hot reload
├── go.work             # Go workspace configuration
├── go.mod              # Go module definition
├── magefile.go         # Main mage file (currently unused)
└── tailwind.config.js  # Tailwind CSS configuration
```

## Getting Started

### 1. Clone the repository

```bash
git clone https://github.com/Damianko135/playground-go
cd playground-go
```

### 2. Install dependencies

```bash
# Install Go dependencies
go mod download

# Note: Tailwind CSS will be run via npx (no separate installation needed)
```

### 3. Run the development server

```bash
mage dev
```

This command will:
- Generate Templ templates
- Start Tailwind CSS watcher
- Start the development server with hot reload
- Make the application available at http://localhost:3000

## Go Workspace Configuration

This project uses Go workspaces to manage multiple modules:

- **Main module** (`.`) - Contains the main application code
- **Magefiles module** (`./magefiles`) - Contains build automation scripts

The `go.work` file is configured to include both modules:

```go
go 1.24.4

use .
use ./magefiles
```

**Important:** Both modules must be included in the workspace for the build system to work correctly.

## Available Mage Commands

Run `mage -l` to see all available commands:

```bash
# Development
mage dev          # Start development server with hot reload and Tailwind
mage watch        # Watch templ files for changes and regenerate
mage hot          # Run server with Air hot reload only

# Building & Testing
mage build        # Build the application binary
mage run          # Run the built binary
mage test         # Run all Go tests
mage ci           # Run complete CI pipeline (lint, test, build)

# Code Quality
mage lint         # Run golangci-lint for code quality checks
mage fmt          # Format all Go code using gofmt

# Utilities
mage clean        # Remove build artifacts
mage cleanDev     # Clean development artifacts and caches
mage tools        # Install all required development tools
mage install      # Install binary to GOBIN or /usr/local/bin

# Documentation
mage docs:gen     # Generate Markdown documentation from Go code
```

## Development

### Hot Reload

The development server includes hot reload functionality powered by Air. Any changes to Go files will automatically rebuild and restart the server.

### Templating

This project uses [Templ](https://templ.guide/) for type-safe HTML templating. Template files are located in the `views/` directory and have the `.templ` extension.

To generate Go code from templates:
```bash
templ generate
```

### Styling

Tailwind CSS is used for styling. The main CSS file is located at `internal/app.css` and is compiled to `static/index.css`.

## Troubleshooting

### "module not in workspace" Error

If you encounter an error like "current directory is contained in a module that is not one of the workspace modules", ensure your `go.work` file includes both modules:

```go
go 1.24.4

use .
use ./magefiles
```

### Port Already in Use

If port 3000 is already in use, you can modify the port in the server configuration or stop the process using that port.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests and ensure everything works
5. Submit a pull request

## License

[Add your license information here]