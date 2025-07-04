# CI/CD & Development Workflows

Local development integration with CI pipeline tasks for consistent development and testing.

## 🎯 Overview

The CI/CD system ensures **the same tasks run both locally and in CI**:
- ✅ **Reproducible builds** - What works locally works in CI
- ✅ **Faster feedback** - Catch issues before pushing
- ✅ **Consistent environments** - Same tools, same versions, same results

## 🚀 Quick Start

```bash
# Run quick checks before committing
task ci:quick

# Run all PR checks before creating a PR
task ci:pr

# Run complete CI pipeline locally
task ci:full

# Show CI status dashboard
task ci:status
```

## 📋 CI Commands Reference

### Quick Development Commands

| Command | Description | Use Case | Duration |
|---------|-------------|----------|----------|
| `task ci:quick` | Test + Lint only | Quick pre-commit check | 2-3 min |
| `task ci:pre-commit` | Pre-commit hook equivalent | Git hook integration | 2-3 min |
| `task ci:pr` | All PR checks | Before creating PR | 5-10 min |
| `task ci:main` | All main branch checks | Before merging to main | 8-12 min |
| `task ci:full` | Complete CI simulation | Full pipeline test | 10-15 min |

### Individual Job Commands

| Command | Description | CI Equivalent |
|---------|-------------|---------------|
| `task ci:test` | Run test job | GitHub Actions test job |
| `task ci:lint` | Run lint job | GitHub Actions lint job |
| `task ci:build` | Build for current platform | GitHub Actions build job |
| `task ci:build:all` | Build for all platforms | GitHub Actions build matrix |
| `task ci:security` | Run security scans | GitHub Actions security job |
| `task ci:docker` | Docker build + scan | GitHub Actions docker job |
| `task ci:codeql` | CodeQL analysis | GitHub Actions codeql job |

## 🐳 Docker Workflows

### Basic Operations
```bash
# Build image (single platform)
task docker:build

# Build multi-platform (like CI)
task docker:build:multi

# Run container locally
task docker:run

# Test image functionality
task docker:test
```

### Security & Registry
```bash
# Scan for vulnerabilities
task docker:scan

# Scan with SARIF output (like CI)
task docker:scan:sarif

# Complete Docker CI workflow
task docker:ci

# Login to GitHub Container Registry
task docker:login

# Push to registry
task docker:push
```

## 🔒 Security Analysis

### CodeQL Setup & Analysis
```bash
# Install CodeQL CLI
task codeql:install

# Run complete CodeQL analysis (like CI)
task codeql:ci

# Run security analysis only
task codeql:analyze:security

# Run quality analysis only
task codeql:analyze:quality
```

### Results Management
```bash
# View results in human-readable format
task codeql:view:results

# Show CodeQL information
task codeql:info

# Clean up CodeQL files
task codeql:clean
```

## 🔄 Workflow Examples

### Daily Development

```bash
# Before starting work
task setup

# During development (with hot reload)
task dev:hot

# Quick check before committing
task ci:quick

# More thorough check
task ci:pr
```

### Release Preparation

```bash
# Full validation before release
task ci:full

# Build all platform binaries
task ci:build:all

# Security scan everything
task ci:security
task docker:scan
```

### Docker Development

```bash
# Build and test Docker image
task docker:build
task docker:test
task docker:scan

# Or run complete Docker workflow
task docker:ci
```

## 🗺️ CI/CD Pipeline Mapping

| Local Command | GitHub Actions Workflow | Description |
|---------------|-------------------------|-------------|
| `task ci:test` | `.github/workflows/ci.yml` (test job) | Unit tests + race detection |
| `task ci:lint` | `.github/workflows/ci.yml` (lint job) | Quality checks + vulnerabilities |
| `task ci:build` | `.github/workflows/ci.yml` (build job) | Multi-platform builds |
| `task ci:security` | `.github/workflows/ci.yml` (security job) | Security scans |
| `task docker:ci` | `.github/workflows/docker.yml` | Docker build + scan |
| `task codeql:ci` | `.github/workflows/codeql.yml` | CodeQL analysis |

## 🛠️ Tool Management

All tools are automatically installed when needed:

| Tool | Installation Method | Purpose |
|------|-------------------|---------|
| **golangci-lint** | Auto-installed via task | Code linting |
| **govulncheck** | Auto-installed via task | Vulnerability scanning |
| **gosec** | Auto-installed via task | Security analysis |
| **Trivy** | Homebrew (macOS) / apt (Linux) | Container scanning |
| **CodeQL CLI** | GitHub releases | Code analysis |
| **gotestsum** | Auto-installed via task | Test output formatting |

## 📁 Configuration Sync

| File | Purpose | Local/CI Sync |
|------|---------|---------------|
| `.golangci.yml` | Linting rules | ✅ Same config |
| `Dockerfile` | Container build | ✅ Same build |
| `go.mod` | Go dependencies | ✅ Same versions |
| `Taskfile.yml` | Task definitions | ✅ Same commands |
| `versions.yml` | Tool versions | ✅ Auto-synced |

## 🚨 Troubleshooting

### Common Issues

1. **Tool not found**
   - Tasks auto-install tools
   - Restart shell after first run
   - Check with `task shared:info:environment`

2. **Permission denied**
   - Ensure Docker is running
   - Check user permissions for Docker
   - Use `sudo` if needed for system tools

3. **Network issues**
   - Tools need internet for downloads
   - Check firewall/proxy settings
   - Some tools cache downloads locally

4. **Version mismatches**
   - Run `task shared:version:check`
   - Sync with `task shared:version:sync`
   - Check `versions.yml` for correct versions

### Getting Help

```bash
# Show environment info
task shared:info:environment

# Show all available tasks
task shared:info:tasks

# Show CI status
task ci:status
```

## 🎯 Best Practices

1. **Pre-commit**: Always run `task ci:quick` before committing
2. **Pre-PR**: Run `task ci:pr` before creating pull requests
3. **Docker testing**: Test containers locally with `task docker:ci`
4. **Security scanning**: Include security checks in regular workflow
5. **Version management**: Keep versions updated via `versions.yml`
6. **Tool consistency**: Let the system manage tool installations 