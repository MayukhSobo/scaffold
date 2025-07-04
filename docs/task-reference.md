# Task Reference

Complete reference of all available tasks organized by category.

## 🚀 Quick Reference

```bash
# Show all available tasks
task

# Show detailed help
task help:help

# Setup complete environment
task setup
```

## 📦 Build Tasks

| Command | Description |
|---------|-------------|
| `task build:debug` | Build a development binary with debug symbols and race detection |
| `task build:release:linux` | Build an optimized release binary for Linux |
| `task build:release:darwin` | Build an optimized release binary for macOS |
| `task build:release:windows` | Build an optimized release binary for Windows |
| `task build:release:all` | Build release binaries for all platforms |

## 🏃 Development Tasks

| Command | Description |
|---------|-------------|
| `task dev:run` | Run the application using `configs/local.yml` |
| `task dev:run:prod` | Run the application using `configs/prod.yml` |
| `task dev:hot` | Run with **hot-reloading** using `air` |

## 🧪 Test Tasks

| Command | Description |
|---------|-------------|
| `task test:all` | Run all tests using `gotestsum` |
| `task test:unit` | Run only unit tests |
| `task test:integration` | Run only integration tests |
| `task test:race` | Run tests with the race detector enabled |
| `task test:benchmark` | Run benchmark tests |
| `task test:coverage` | Generate a standard HTML coverage report |
| `task test:coverage:open` | Generate comprehensive coverage reports and open in browser |

## 🔍 Code Quality Tasks

| Command | Description |
|---------|-------------|
| `task quality:lint` | Run `golangci-lint` to find code issues |
| `task quality:fix` | Run linter with auto-fix enabled |
| `task quality:fmt` | Format all Go source files with `gofmt` |
| `task quality:vet` | Run `go vet` to analyze source code |
| `task quality:all` | Run all quality checks (`fmt`, `vet`, `lint`) |
| `task quality:gosec` | Run security analysis with `gosec` |

## 🧹 Cleanup Tasks

| Command | Description |
|---------|-------------|
| `task clean:all` | Clean all build artifacts and caches |
| `task clean:debug` | Clean development build artifacts only |
| `task clean:release:linux` | Clean Linux release build artifacts only |
| `task clean:release:darwin` | Clean macOS release build artifacts only |
| `task clean:release:windows` | Clean Windows release build artifacts only |
| `task clean:release:all` | Clean all release build artifacts |

## 📦 Dependency Management

| Command | Description |
|---------|-------------|
| `task deps:install` | Download and tidy Go module dependencies |
| `task deps:update` | Update all dependencies to the latest versions |
| `task deps:patch` | Update patch versions only (safer) |
| `task deps:verify` | Verify dependencies and checksums |
| `task deps:vulnerabilities` | Check for known security vulnerabilities |
| `task deps:outdated` | Check for outdated dependencies |
| `task deps:licenses` | Generate license report |

## 🐳 Docker Tasks

| Command | Description |
|---------|-------------|
| `task docker:build` | Build a production-ready Docker image |
| `task docker:build:multi` | Build multi-platform Docker image |
| `task docker:build:ci` | Build with CI-style caching |
| `task docker:run` | Run the application in a Docker container |
| `task docker:test` | Test the Docker image functionality |
| `task docker:scan` | Scan Docker image for vulnerabilities |
| `task docker:scan:sarif` | Scan with SARIF output format |
| `task docker:ci` | Complete Docker CI workflow |
| `task docker:login` | Login to GitHub Container Registry |
| `task docker:push` | Push image to registry |

## 🔒 Security & Analysis

| Command | Description |
|---------|-------------|
| `task codeql:install` | Install CodeQL CLI tool |
| `task codeql:ci` | Run complete CodeQL analysis |
| `task codeql:analyze:security` | Run security analysis only |
| `task codeql:analyze:quality` | Run quality analysis only |
| `task codeql:view:results` | View CodeQL results |
| `task codeql:info` | Show CodeQL information |
| `task codeql:clean` | Clean CodeQL artifacts |

## 🔄 CI/CD Tasks

| Command | Description |
|---------|-------------|
| `task ci:quick` | Run quick checks (test + lint) |
| `task ci:pr` | Run all PR checks |
| `task ci:main` | Run all main branch checks |
| `task ci:full` | Run complete CI pipeline |
| `task ci:test` | Run CI test job |
| `task ci:lint` | Run CI lint job |
| `task ci:build` | Run CI build job |
| `task ci:build:all` | Build for all platforms |
| `task ci:security` | Run CI security job |
| `task ci:docker` | Run CI Docker job |
| `task ci:status` | Show CI status dashboard |

## 🔧 Setup & Environment

| Command | Description |
|---------|-------------|
| `task setup` | Complete environment setup |
| `task shared:setup` | Setup development environment |
| `task shared:setup:tools` | Install development tools |
| `task shared:setup:python` | Setup Python virtual environment |
| `task shared:setup:python:activate` | Show Python venv activation instructions |
| `task shared:validate:go` | Validate Go environment |
| `task shared:validate:docker` | Validate Docker environment |
| `task shared:validate:python` | Validate Python environment |

## 🔢 Version Management

| Command | Description |
|---------|-------------|
| `task shared:version:show` | Show all versions with beautiful formatting |
| `task shared:version:sync` | Sync versions across all files |
| `task shared:version:sync:dry-run` | Preview version sync changes |
| `task shared:version:check` | Check version consistency |

## 🧹 Cleanup & Maintenance

| Command | Description |
|---------|-------------|
| `task shared:cleanup:go` | Clean Go caches |
| `task shared:cleanup:build` | Clean build artifacts |
| `task shared:cleanup:docker` | Clean Docker system |
| `task shared:cleanup:reports` | Clean all reports |
| `task shared:cleanup:logs` | Clean log files |

## ℹ️ Information & Help

| Command | Description |
|---------|-------------|
| `task shared:info:environment` | Show environment information |
| `task shared:info:tasks` | Show detailed task reference |
| `task shared:info:go` | Show Go environment info |
| `task shared:info:docker` | Show Docker environment info |
| `task shared:info:python` | Show Python environment info |
| `task help:help` | Show detailed help for all tasks |

## ⚙️ Configuration & Validation

| Command | Description |
|---------|-------------|
| `task config:validate` | Validate all `.yml` files in the `configs` directory |
| `task config:validate:local` | Validate local configuration |
| `task config:validate:prod` | Validate production configuration |

## 🎯 Usage Examples

### Daily Development Workflow

```bash
# Setup environment (first time)
task setup

# Start development with hot reload
task dev:hot

# Quick checks before committing
task ci:quick

# More thorough checks before PR
task ci:pr
```

### Release Preparation

```bash
# Run complete CI pipeline
task ci:full

# Build for all platforms
task ci:build:all

# Run security analysis
task ci:security

# Build and scan Docker image
task docker:ci
```

### Maintenance

```bash
# Update dependencies
task deps:update

# Check for vulnerabilities
task deps:vulnerabilities

# Sync versions
task shared:version:sync

# Clean up everything
task clean:all
```

## 📝 Task Arguments

Some tasks accept arguments:

```bash
# Install specific tools
task shared:setup:tools -- gotestsum golangci-lint

# Run specific tests
task test:unit -- -v -run TestUserService

# Build with custom flags
task build:debug -- -ldflags="-X main.version=dev"
```

## 🔧 Advanced Usage

### Environment Variables

Set these environment variables to customize behavior:

```bash
# Skip tool installation checks
export SKIP_TOOL_CHECK=true

# Force tool reinstallation
export FORCE_TOOL_INSTALL=true

# Custom Go version
export GO_VERSION=1.22.0

# Custom Docker tag
export DOCKER_TAG=custom-tag
```

### Task Options

```bash
# Run with verbose output
task --verbose ci:pr

# Run in specific directory
task --dir /path/to/project ci:pr

# Show what would be executed
task --dry-run ci:pr
```

## 🚨 Troubleshooting

### Common Task Issues

1. **Task not found**: Run `task --list` to see available tasks
2. **Permission denied**: Check file permissions and Docker access
3. **Tool not installed**: Run `task shared:setup:tools` to install all tools
4. **Version mismatch**: Run `task shared:version:check` and `task shared:version:sync`

### Getting Help

```bash
# Show all tasks
task

# Show help for specific task
task help:help

# Show environment info
task shared:info:environment

# Show task details
task shared:info:tasks
``` 