# Version Management System

Centralized version management with automatic synchronization across all project files.

## 🔢 Overview

All tool versions are managed in `versions.yml` and automatically synced to:
- GitHub Actions workflows
- Dockerfile
- go.mod
- Configuration files

## 🚀 Quick Commands

```bash
# Show all versions in organized format
task shared:version:show

# Sync versions across all files
task shared:version:sync

# Preview changes before syncing
task shared:version:sync:dry-run

# Check version consistency
task shared:version:check
```

## 📁 Version File Structure

The `versions.yml` file uses a structured format:

```yaml
# Core language versions
go: "1.24.4"
node: "20"
python: "3.12"

# Development tools
tools:
  golangci-lint: "v2.2.1"
  gotestsum: "latest"
  gosec: "latest"
  govulncheck: "latest"
  air: "latest"
  gocov: "latest"
  gocov-html: "latest"
  go-cover-treemap: "latest"
  trivy: "latest"

# GitHub Actions
actions:
  setup-go: "v5"
  setup-task: "v3"
  setup-docker-buildx: "v3"
  docker-metadata: "v5"
  docker-build-push: "v6"
  github-script: "v7"
  upload-artifact: "v4"
  download-artifact: "v4"
  cache: "v4"

# Security tools
security:
  codeql-cli: "2.16.1"
  trivy: "latest"
  gosec: "latest"
  govulncheck: "latest"

# Build tools
build:
  docker-buildx: "latest"
  task: "3.x"
  schema_version: "1.0.0"
```

## 🔄 How It Works

### 1. Python-Based Version Helper

The system uses a robust Python script (`scripts/version-helper.py`) that:
- Parses YAML with PyYAML for reliability
- Provides fallback parsing if PyYAML is unavailable
- Offers a clean command-line interface
- Integrates seamlessly with Task

### 2. Automatic Synchronization

The sync script (`scripts/sync-versions.py`) automatically updates:

**GitHub Workflows:**
- `go-version` in setup-go actions
- Environment variables like `GO_VERSION`
- Action versions (`uses: actions/setup-go@v5`)

**Configuration Files:**
- `Dockerfile` - Go base image version
- `go.mod` - Go version requirement
- `Taskfile.yml` - Version variables

**Scripts:**
- All installation scripts use centralized versions
- Tool installation respects version specifications

### 3. Consistency Checking

The system can detect version inconsistencies across files:
- Compares expected vs actual versions
- Shows detailed diff with rich formatting
- Prevents version drift between files

## 📖 Usage Examples

### Viewing Versions

```bash
# Show all versions in organized format
task shared:version:show

# Get specific version programmatically
python scripts/run.py version-helper get go
python scripts/run.py version-helper get tools.golangci-lint
```

### Updating Versions

1. **Edit `versions.yml`** - Update the desired version
2. **Preview changes:**
   ```bash
   task shared:version:sync:dry-run
   ```
3. **Apply changes:**
   ```bash
   task shared:version:sync
   ```

### Checking Consistency

```bash
# Verify all files use correct versions
task shared:version:check
```

## 🛡️ Safety Features

- **Dry-run mode**: Preview changes before applying
- **Consistency checking**: Detect version mismatches
- **Backup-free operation**: Git tracks all changes
- **Verbose output**: See exactly what's being changed
- **Error handling**: Clear error messages for issues

## 🔧 Integration

### In Task Files

```yaml
vars:
  GO_VERSION:
    sh: python scripts/run.py version-helper get go
  GOLANGCI_LINT_VERSION:
    sh: python scripts/run.py version-helper get tools.golangci-lint
```

### In Scripts

```bash
# Get versions in shell scripts
GO_VERSION=$(python scripts/run.py version-helper get go)
LINT_VERSION=$(python scripts/run.py version-helper get tools.golangci-lint)
```

### In GitHub Workflows

Versions are automatically synced, so workflows always use the correct versions from `versions.yml`.

## 🎯 Benefits

- **Single source of truth**: All versions defined in one place
- **Automatic synchronization**: No manual updates across files
- **Consistency guarantee**: Prevents version drift
- **Clear output**: Organized display of version information
- **Error prevention**: Catches inconsistencies before they cause issues
- **CI/CD integration**: Workflows always use correct versions 