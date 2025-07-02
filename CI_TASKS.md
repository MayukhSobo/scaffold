# CI Tasks - Local Development Integration

This document explains how to run CI pipeline tasks locally for consistent development and testing.

## Overview

We've streamlined the CI/CD pipeline so that **the same tasks run both locally and in CI**. This ensures:
- ✅ **Reproducible builds** - What works locally works in CI
- ✅ **Faster feedback** - Catch issues before pushing
- ✅ **Consistent environments** - Same tools, same versions, same results

## Quick Start

```bash
# Run quick checks before committing
task ci:quick

# Run all PR checks before creating a PR
task ci:pr

# Run complete CI pipeline locally
task ci:full

# Show all available CI commands
task ci:status
```

## Available CI Commands

### Quick Development Commands

| Command | Description | Use Case |
|---------|-------------|----------|
| `task ci:quick` | Test + Lint only | Quick pre-commit check |
| `task ci:pre-commit` | Pre-commit hook equivalent | Git hook integration |
| `task ci:pr` | All PR checks | Before creating PR |
| `task ci:main` | All main branch checks | Before merging to main |
| `task ci:full` | Complete CI simulation | Full pipeline test |

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

## Docker Tasks

### Basic Docker Operations
```bash
# Build image (single platform)
task docker:build

# Build multi-platform (like CI)
task docker:build:multi

# Run container
task docker:run

# Test image (like CI)
task docker:test
```

### Security Scanning
```bash
# Scan for vulnerabilities
task docker:scan

# Scan with SARIF output (like CI)
task docker:scan:sarif

# Complete Docker CI workflow
task docker:ci
```

### Registry Operations
```bash
# Login to GitHub Container Registry
task docker:login

# Push to registry
task docker:push
```

## CodeQL Analysis

### Setup and Analysis
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

### Results and Management
```bash
# View results in human-readable format
task codeql:view:results

# Show CodeQL info
task codeql:info

# Clean up CodeQL files
task codeql:clean
```

## Workflow Examples

### Before Committing
```bash
# Quick check (2-3 minutes)
task ci:quick
```

### Before Creating PR
```bash
# Complete PR validation (5-10 minutes)
task ci:pr
```

### Before Merging to Main
```bash
# Full pipeline simulation (10-15 minutes)
task ci:full
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

### Security Analysis
```bash
# Run all security checks
task ci:security

# Run CodeQL analysis
task codeql:ci

# Scan Docker image
task docker:scan
```

## CI/CD Pipeline Mapping

| Local Command | GitHub Actions Workflow | Description |
|---------------|-------------------------|-------------|
| `task ci:test` | `.github/workflows/ci.yml` (test job) | Unit tests + race detection |
| `task ci:lint` | `.github/workflows/ci.yml` (lint job) | Quality checks + vulnerabilities |
| `task ci:build` | `.github/workflows/ci.yml` (build job) | Multi-platform builds |
| `task ci:security` | `.github/workflows/ci.yml` (security job) | Security scans |
| `task docker:ci` | `.github/workflows/docker.yml` | Docker build + scan |
| `task codeql:ci` | `.github/workflows/codeql.yml` | CodeQL analysis |

## Tool Installation

All tools are automatically installed when needed:

- **golangci-lint**: Auto-installed via task
- **govulncheck**: Auto-installed via task  
- **gosec**: Auto-installed via task
- **Trivy**: Auto-installed via Homebrew (macOS) or apt (Linux)
- **CodeQL CLI**: Auto-downloaded from GitHub releases

## Configuration Files

| File | Purpose | Local/CI Sync |
|------|---------|---------------|
| `.golangci.yml` | Linting rules | ✅ Same config |
| `Dockerfile` | Container build | ✅ Same build |
| `go.mod` | Go dependencies | ✅ Same versions |
| `Taskfile.yml` | Task definitions | ✅ Same commands |

## Troubleshooting

### Common Issues

1. **Tool not found**: Tasks auto-install tools, but you might need to restart your shell
2. **Permission denied**: Make sure Docker is running and you have permissions
3. **Network issues**: Some tools need internet access for downloads/updates

### Getting Help

```bash
# Show all available tasks
task --list

# Show CI status and commands
task ci:status

# Show CodeQL info
task codeql:info

# Clean up everything and start fresh
task ci:clean
```

## Benefits

### For Developers
- 🚀 **Faster feedback**: Catch issues locally before CI
- 🔄 **Consistent experience**: Same tools, same results
- 🛠️ **Easy debugging**: Run individual CI jobs locally
- ⚡ **Quick iterations**: No need to push to test

### For Teams
- 📊 **Reduced CI failures**: Issues caught locally first
- 🔧 **Easier onboarding**: Clear, documented workflows
- 🎯 **Standardized processes**: Everyone uses same tools
- 💰 **Lower CI costs**: Fewer failed CI runs

---

**Pro Tip**: Add `task ci:pre-commit` to your Git pre-commit hooks for automatic validation! 