# Development Guide

Complete guide for development workflows, testing, and building.

## 🚀 Getting Started

### Prerequisites

- [Go](https://golang.org/dl/) (version 1.22+ recommended)
- [Task](https://taskfile.dev/installation/)
- [Docker](https://www.docker.com/get-started) (for containerized builds)
- [Python 3](https://www.python.org/downloads/) (for enhanced tooling)

### First-Time Setup

```bash
# Clone the repository
git clone https://github.com/thedatageek/scaffold.git
cd scaffold

# Setup complete development environment
task setup
```

This will:
1. **Install Go dependencies** (`go mod download`)
2. **Install development tools** (golangci-lint, gotestsum, etc.)
3. **Setup Python virtual environment** with enhanced scripts
4. **Validate environment** setup

## 📊 Daily Development Workflow

### Starting Development

```bash
# Start development with hot reload
task dev:hot

# Or run normally
task dev:run

# Run with production config
task dev:run:prod
```

### Before Committing

```bash
# Quick checks (2-3 minutes)
task ci:quick

# More thorough checks (5-10 minutes)  
task ci:pr
```

### Testing

```bash
# Run all tests
task test:all

# Run specific test types
task test:unit
task test:integration
task test:race

# Generate and view coverage
task test:coverage:open
```

### Code Quality

```bash
# Format code
task quality:fmt

# Run linter
task quality:lint

# Auto-fix issues
task quality:fix

# Run all quality checks
task quality:all
```

## 🏗️ Building

### Development Builds

```bash
# Debug build with race detection
task build:debug

# Run the debug binary
./build/scaffold
```

### Release Builds

```bash
# Build for current platform
task build:release:linux    # Linux
task build:release:darwin   # macOS  
task build:release:windows  # Windows

# Build for all platforms
task build:release:all
```

Built binaries are located in:
- `build/linux/scaffold-amd64-linux`
- `build/darwin/scaffold-amd64-darwin`
- `build/windows/scaffold-amd64-windows.exe`

## 🐳 Docker Development

### Local Docker Builds

```bash
# Build Docker image
task docker:build

# Run container
task docker:run

# Test the image
task docker:test
```

### Multi-Platform Builds

```bash
# Build for multiple architectures (like CI)
task docker:build:multi

# Build with CI-style caching
task docker:build:ci
```

### Security Scanning

```bash
# Scan for vulnerabilities
task docker:scan

# Scan with SARIF output
task docker:scan:sarif

# Complete Docker CI workflow
task docker:ci
```

## 🧪 Testing & Coverage

### Test Types

| Command | Description | When to Use |
|---------|-------------|-------------|
| `task test:unit` | Unit tests only | Fast feedback during development |
| `task test:integration` | Integration tests only | Testing component interactions |
| `task test:race` | Race condition detection | Before committing concurrent code |
| `task test:benchmark` | Performance benchmarks | Performance optimization |

### Coverage Analysis

```bash
# Generate basic coverage
task test:coverage

# Interactive coverage with multiple views
task test:coverage:open
```

This provides:
- **Standard HTML reports**
- **Function-level reports** via `gocov`
- **Visual treemaps** via `go-cover-treemap`

## 🔍 Code Quality & Security

### Linting

```bash
# Run linter with current config
task quality:lint

# Auto-fix issues where possible
task quality:fix

# Check specific aspects
task quality:vet    # go vet
task quality:fmt    # formatting check
```

### Security Analysis

```bash
# Run security scanner
task quality:gosec

# Check for vulnerabilities
task deps:vulnerabilities

# Run complete security suite
task ci:security
```

### CodeQL Analysis

```bash
# Install CodeQL CLI
task codeql:install

# Run complete analysis
task codeql:ci

# Run specific analysis types
task codeql:analyze:security
task codeql:analyze:quality
```

## 📦 Dependency Management

### Managing Dependencies

```bash
# Install/update dependencies
task deps:install

# Update all dependencies
task deps:update

# Update patch versions only (safer)
task deps:patch

# Verify dependencies
task deps:verify
```

### Security & Licensing

```bash
# Check for vulnerabilities
task deps:vulnerabilities

# Check for outdated packages
task deps:outdated

# Generate license report
task deps:licenses
```

## 🧹 Maintenance & Cleanup

### Regular Cleanup

```bash
# Clean all build artifacts
task clean:all

# Clean specific artifacts
task clean:debug
task clean:release:all
```

### Deep Cleanup

```bash
# Clean Go caches
task shared:cleanup:go

# Clean Docker system
task shared:cleanup:docker

# Clean all reports
task shared:cleanup:reports
```

## 🔧 Configuration

### Application Configuration

Configuration files are in `configs/`:
- `local.yml` - Local development settings
- `prod.yml` - Production settings

### Development Tools Configuration

| File | Purpose |
|------|---------|
| `.golangci.yml` | Linter configuration |
| `.air.toml` | Hot-reload configuration |
| `Taskfile.yml` | Task definitions |
| `versions.yml` | Tool versions |

## 🚨 Troubleshooting

### Common Issues

1. **Tools not found**:
   ```bash
   task shared:setup:tools
   # Restart shell after installation
   ```

2. **Python venv issues**:
   ```bash
   task python
   source .venv/bin/activate
   ```

3. **Docker permission issues**:
   ```bash
   # Ensure Docker is running
   docker version
   ```

4. **Version inconsistencies**:
   ```bash
   task shared:version:check
   task shared:version:sync
   ```

### Getting Help

```bash
# Show all available tasks
task

# Show environment information
task shared:info:environment

# Show detailed task reference
task shared:info:tasks
```

## 🎯 Best Practices

### Development
1. **Use hot-reload** during development (`task dev:hot`)
2. **Run quick checks** before each commit (`task ci:quick`)
3. **Use coverage analysis** to guide testing (`task test:coverage:open`)
4. **Keep dependencies updated** (`task deps:patch`)

### Code Quality
1. **Auto-format code** regularly (`task quality:fmt`)
2. **Fix linter issues** promptly (`task quality:fix`)
3. **Run security scans** periodically (`task quality:gosec`)
4. **Check for vulnerabilities** in dependencies

### CI/CD
1. **Test locally** before pushing (`task ci:pr`)
2. **Use same commands** as CI pipeline
3. **Keep tool versions** synchronized
4. **Test Docker builds** locally (`task docker:ci`)

### Documentation
1. **Update documentation** when changing workflows
2. **Document new tasks** with clear descriptions
3. **Keep examples** up to date
4. **Test documentation** examples 