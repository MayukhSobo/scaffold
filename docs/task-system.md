# Task System & Automation

This project uses a unified task system that ensures consistency between local development and CI/CD workflows. The system features centralized tool management, version control, and enhanced Python scripts.

## 🐍 Python Scripts

All development scripts are written in Python with enhanced features:
- **Progress tracking**: Shows installation progress and status
- **Auto venv activation**: Scripts automatically use the Python virtual environment
- **Enhanced error handling**: Clear error messages and recovery suggestions
- **Cross-platform support**: Works on macOS, Linux, and Windows

## 🚀 Quick Start

```bash
# Setup complete development environment (includes Python venv)
task setup

# Show all available tasks with beautiful formatting
task

# Run quick checks before committing
task ci:quick

# Run all PR checks
task ci:pr

# Install specific development tools
task shared:setup:tools -- gotestsum golangci-lint
```

## 📋 Essential Commands

| Command | Description | Use Case |
|---------|-------------|----------|
| `task setup` | Complete environment setup | First-time setup |
| `task ci:quick` | Test + Lint only | Quick pre-commit check |
| `task ci:pr` | All PR checks | Before creating PR |
| `task test:coverage:open` | Interactive coverage | View detailed coverage |
| `task quality:all` | Format + lint + fix | Code quality |
| `task docker:ci` | Docker build + scan | Container workflow |

## 🛠️ Development Tools

The system automatically installs and manages:
- **golangci-lint**: Code linting with comprehensive rules
- **gotestsum**: Beautiful test output formatting
- **gosec**: Security vulnerability scanning
- **govulncheck**: Go vulnerability database scanning
- **air**: Live-reloading for development
- **trivy**: Container security scanning
- **gocov tools**: Enhanced coverage reporting

## 🐳 Docker Integration

```bash
# Build and test Docker image
task docker:build
task docker:test

# Security scanning
task docker:scan

# Complete Docker CI workflow
task docker:ci
```

## Task Categories

### Setup & Environment
- `shared:setup` - Complete environment setup
- `shared:setup:tools` - Install specific tools
- `shared:setup:python` - Setup Python virtual environment
- `shared:validate:go` - Validate Go environment
- `shared:validate:docker` - Validate Docker environment

### Build & Development
- `build:debug` - Build debug binary with race detection
- `build:release:*` - Build optimized binaries for different platforms
- `dev:run` - Run application with local config
- `dev:hot` - Run with hot reload

### Testing
- `test:all` - Run all tests
- `test:unit` - Run unit tests only
- `test:integration` - Run integration tests only
- `test:race` - Run tests with race detection
- `test:coverage` - Generate coverage report
- `test:coverage:open` - Generate enhanced coverage and serve

### Code Quality
- `quality:lint` - Run linter
- `quality:fix` - Run linter with auto-fix
- `quality:fmt` - Format code
- `quality:vet` - Run go vet
- `quality:all` - Run all quality checks with fixes

### CI/CD Simulation
- `ci:test` - Run CI test job
- `ci:lint` - Run CI lint job
- `ci:build` - Run CI build job
- `ci:security` - Run CI security job
- `ci:pr` - Run all PR checks
- `ci:main` - Run all main branch checks
- `ci:quick` - Run quick checks (test + lint)
- `ci:full` - Run complete CI pipeline

### Cleanup
- `clean:all` - Clean all artifacts
- `shared:cleanup:go` - Clean Go caches
- `shared:cleanup:build` - Clean build artifacts
- `shared:cleanup:docker` - Clean Docker system

## Getting Help

```bash
# Show all available tasks
task

# Show environment information
task shared:info:environment

# Show detailed task reference
task shared:info:tasks
``` 